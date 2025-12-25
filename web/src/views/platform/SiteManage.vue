<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>站点管理</span>
          <el-button type="primary" :icon="Plus" @click="handleOpenDialog()"> 新建站点 </el-button>
        </div>
      </template>

      <el-table v-loading="loading" :data="schoolList" border stripe style="width: 100%">
        <el-table-column prop="ID" label="ID" width="80" align="center" />
        <el-table-column prop="Name" label="站点名称" min-width="150" />
        <el-table-column prop="ContactName" label="负责人" min-width="120" />
        <el-table-column prop="ContactPhone" label="联系电话" min-width="150" />
        <el-table-column prop="Address" label="地址" min-width="200" show-overflow-tooltip />
        <el-table-column prop="AdminUsername" label="管理员账号" min-width="150" />
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch v-model="row.IsEnabled" @change="handleStatusChange(row)" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" :icon="Edit" link @click="handleEdit(row)">编辑</el-button>
            <el-button
              type="warning"
              :icon="RefreshRight"
              link
              @click="handleOpenResetPwdDialog(row.AdminUserID)"
              >重置管理员密码</el-button
            >
            <el-button type="danger" :icon="Delete" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <Pagination
          :total="total"
          v-model:page="pagination.page"
          v-model:limit="pagination.pageSize"
          @pagination="getList"
        />
      </div>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="dialogMode === 'create' ? '新建站点' : '编辑站点'"
      width="600px"
      @close="handleCloseDialog"
    >
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px">
        <el-form-item label="站点名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入学校名称" />
        </el-form-item>
        <el-form-item label="负责人" prop="contactName">
          <el-input v-model="formData.contactName" placeholder="请输入负责人姓名" />
        </el-form-item>
        <el-form-item label="联系电话" prop="contactPhone">
          <el-input v-model="formData.contactPhone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="地址" prop="address">
          <el-input v-model="formData.address" type="textarea" placeholder="请输入地址" />
        </el-form-item>
        <template v-if="dialogMode === 'edit'">
          <el-form-item label="状态" prop="isEnabled">
            <el-radio-group v-model="formData.isEnabled">
              <el-radio :label="true">启用</el-radio>
              <el-radio :label="false">禁用</el-radio>
            </el-radio-group>
          </el-form-item>
        </template>
        <el-divider v-if="dialogMode === 'create'">管理员信息</el-divider>
        <template v-if="dialogMode === 'create'">
          <el-form-item label="管理员账号" prop="adminUsername">
            <el-input v-model="formData.adminUsername" placeholder="用于登录学校后台的管理员账号" />
          </el-form-item>
          <el-form-item label="管理员姓名" prop="adminRealName">
            <el-input v-model="formData.adminRealName" placeholder="请输入管理员真实姓名" />
          </el-form-item>
          <el-form-item label="初始密码" prop="adminPassword">
            <el-input
              v-model="formData.adminPassword"
              type="password"
              show-password
              placeholder="请输入初始密码"
            />
          </el-form-item>
          <el-form-item label="确认密码" prop="confirmAdminPassword">
            <el-input
              v-model="formData.confirmAdminPassword"
              type="password"
              show-password
              placeholder="请再次输入密码"
            />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="isSubmitting" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog
      v-model="resetPwdDialogVisible"
      title="重置管理员密码"
      width="500px"
      @close="handleCloseResetPwdDialog"
    >
      <el-form
        ref="resetPwdFormRef"
        :model="resetPwdFormData"
        :rules="resetPwdFormRules"
        label-width="100px"
      >
        <el-form-item label="新密码" prop="password">
          <el-input
            v-model="resetPwdFormData.password"
            type="password"
            show-password
            placeholder="请输入新密码"
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="resetPwdFormData.confirmPassword"
            type="password"
            show-password
            placeholder="请再次输入新密码"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="resetPwdDialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="isResettingPassword" @click="handleResetPwdSubmit">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted } from 'vue';
  import { ElMessage, ElMessageBox } from 'element-plus';
  import { Plus, Edit, Delete, RefreshRight } from '@element-plus/icons-vue';
  import type { FormInstance, FormRules } from 'element-plus';
  import Pagination from '@/components/Pagination/index.vue';
  import {
    listSchoolsApi,
    createSchoolApi,
    updateSchoolApi,
    deleteSchoolApi,
    type School,
    type CreateSchoolPayload,
    type UpdateSchoolPayload,
  } from '@/api/school';
  import { resetPasswordApi, type ResetPasswordPayload } from '@/api/account';

  const isSubmitting = ref(false);
  const isResettingPassword = ref(false);

  // --- 列表和分页 ---
  const loading = ref(false);
  const schoolList = ref<School[]>([]);
  const total = ref(0);
  const pagination = reactive({
    page: 1,
    pageSize: 20,
  });

  const getList = async () => {
    loading.value = true;
    try {
      const res = await listSchoolsApi(pagination);
      schoolList.value = res.list;
      total.value = res.total;
    } catch (error) {
      console.error(error);
    } finally {
      loading.value = false;
    }
  };

  onMounted(() => {
    getList();
  });

  // --- 对话框和表单 ---
  const dialogVisible = ref(false);
  const dialogMode = ref<'create' | 'edit'>('create');
  const currentSchoolId = ref<number | null>(null);
  const formRef = ref<FormInstance | null>(null);

  const getInitialFormData = (): CreateSchoolPayload & {
    confirmAdminPassword: '';
    isEnabled: boolean;
  } => ({
    name: '',
    contactName: '',
    contactPhone: '',
    address: '',
    isEnabled: true,
    adminUsername: '',
    adminPassword: '',
    confirmAdminPassword: '',
    adminRealName: '',
  });

  const formData = reactive(getInitialFormData());

  const validateAdminPass = (_rule: any, value: any, callback: any) => {
    if (value === '') {
      callback(new Error('请再次输入密码'));
    } else if (value !== formData.adminPassword) {
      callback(new Error('两次输入的密码不一致'));
    } else {
      callback();
    }
  };

  const formRules = reactive<FormRules>({
    name: [{ required: true, message: '请输入站点名称', trigger: 'blur' }],
    adminUsername: [{ required: true, message: '请输入管理员账号', trigger: 'blur' }],
    adminRealName: [{ required: true, message: '请输入管理员姓名', trigger: 'blur' }],
    adminPassword: [
      { required: true, message: '请输入初始密码', trigger: 'blur' },
      { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
    ],
    confirmAdminPassword: [{ required: true, validator: validateAdminPass, trigger: 'blur' }],
  });

  const handleOpenDialog = () => {
    dialogMode.value = 'create';
    dialogVisible.value = true;
  };

  const handleEdit = (row: School) => {
    dialogMode.value = 'edit';
    currentSchoolId.value = row.ID;
    formData.name = row.Name;
    formData.contactName = row.ContactName;
    formData.contactPhone = row.ContactPhone;
    formData.address = row.Address;
    formData.isEnabled = row.IsEnabled;
    dialogVisible.value = true;
  };

  const handleCloseDialog = () => {
    formRef.value?.resetFields();
    Object.assign(formData, getInitialFormData());
    currentSchoolId.value = null;
    dialogVisible.value = false;
  };

  const handleSubmit = async () => {
    if (!formRef.value) return;
    await formRef.value.validate(async (valid) => {
      if (valid) {
        isSubmitting.value = true;
        try {
          if (dialogMode.value === 'create') {
            await createSchoolApi(formData);
            ElMessage.success('站点创建成功');
          } else {
            if (!currentSchoolId.value) return;
            const payload: UpdateSchoolPayload = {
              Name: formData.name,
              ContactName: formData.contactName,
              ContactPhone: formData.contactPhone,
              Address: formData.address,
              IsEnabled: formData.isEnabled,
            };
            await updateSchoolApi(currentSchoolId.value, payload);
            ElMessage.success('站点更新成功');
          }
          handleCloseDialog();
          await getList();
        } catch (error) {
          console.error(error);
        } finally {
          isSubmitting.value = false;
        }
      }
    });
  };

  // --- 重置管理员密码对话框 ---
  const resetPwdDialogVisible = ref(false);
  const currentAdminUserId = ref<number | null>(null);
  const resetPwdFormRef = ref<FormInstance | null>(null);
  const resetPwdFormData = reactive({ password: '', confirmPassword: '' });

  const validateResetPass = (_rule: any, value: any, callback: any) => {
    if (value === '') {
      callback(new Error('请再次输入密码'));
    } else if (value !== resetPwdFormData.password) {
      callback(new Error('两次输入的密码不一致'));
    } else {
      callback();
    }
  };

  const resetPwdFormRules = reactive<FormRules>({
    password: [
      { required: true, message: '请输入新密码', trigger: 'blur' },
      { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
    ],
    confirmPassword: [{ required: true, validator: validateResetPass, trigger: 'blur' }],
  });

  const handleOpenResetPwdDialog = (adminUserId: number) => {
    if (!adminUserId) {
      ElMessage.warning('该站点没有关联的管理员账号');
      return;
    }
    currentAdminUserId.value = adminUserId;
    resetPwdDialogVisible.value = true;
  };

  const handleCloseResetPwdDialog = () => {
    resetPwdFormRef.value?.resetFields();
    currentAdminUserId.value = null;
    resetPwdDialogVisible.value = false;
  };

  const handleResetPwdSubmit = async () => {
    if (!resetPwdFormRef.value || !currentAdminUserId.value) return;
    await resetPwdFormRef.value.validate(async (valid) => {
      if (valid) {
        isResettingPassword.value = true;
        try {
          const payload: ResetPasswordPayload = { password: resetPwdFormData.password };
          await resetPasswordApi(currentAdminUserId.value!, payload);
          ElMessage.success('管理员密码重置成功');
          handleCloseResetPwdDialog();
        } catch (error) {
          console.error(error);
        } finally {
          isResettingPassword.value = false;
        }
      }
    });
  };

  // --- 其他操作 ---
  const handleStatusChange = async (row: School) => {
    const actionText = row.IsEnabled ? '启用' : '禁用';
    try {
      await ElMessageBox.confirm(`确定要${actionText} [${row.Name}] 吗？`, '提示', {
        type: 'warning',
      });
      await updateSchoolApi(row.ID, {
        Name: row.Name,
        ContactName: row.ContactName,
        ContactPhone: row.ContactPhone,
        Address: row.Address,
        IsEnabled: row.IsEnabled,
      });
      ElMessage.success(`${actionText}成功`);
    } catch (error) {
      if (error !== 'cancel') ElMessage.error('操作失败');
      await getList(); // 失败或取消时恢复状态
    }
  };

  const handleDelete = async (row: School) => {
    try {
      await ElMessageBox.confirm(`确定要永久删除 [${row.Name}] 吗？此操作不可恢复。`, '危险操作', {
        type: 'error',
      });
      await deleteSchoolApi(row.ID);
      ElMessage.success('删除成功');
      await getList();
    } catch (error) {
      if (error !== 'cancel') ElMessage.error('删除失败');
    }
  };
</script>

<style scoped>
  .app-container {
    height: 100%;
    display: flex;
    flex-direction: column;
  }
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .el-card {
    display: flex;
    flex-direction: column;
    height: 100%;
  }
  :deep(.el-card__body) {
    display: flex;
    flex-direction: column;
    flex-grow: 1;
    overflow: hidden;
  }
  .el-table {
    flex-grow: 1;
    overflow-y: auto;
  }
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
    flex-shrink: 0;
  }
</style>
