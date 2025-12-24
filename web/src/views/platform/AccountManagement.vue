<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>子账号管理</span>
          <el-button type="primary" :icon="Plus" @click="handleOpenDialog()"> 新建账号 </el-button>
        </div>
      </template>

      <el-table v-loading="loading" :data="accountList" border stripe style="width: 100%">
        <el-table-column prop="ID" label="ID" width="80" align="center" />
        <el-table-column prop="Username" label="用户名" min-width="150" />
        <el-table-column prop="RealName" label="真实姓名" min-width="150" />
        <el-table-column prop="Mobile" label="手机号" min-width="150" />
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.Status"
              :active-value="1"
              :inactive-value="2"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="CreatedAt" label="创建时间" min-width="180" />
        <el-table-column label="操作" width="220" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" :icon="Edit" link @click="handleEdit(row)">编辑</el-button>
            <el-button
              type="warning"
              :icon="RefreshRight"
              link
              @click="handleOpenResetPwdDialog(row)"
              >重置密码</el-button
            >
            <el-button type="danger" :icon="Delete" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <Pagination
        :total="total"
        v-model:page="pagination.page"
        v-model:limit="pagination.pageSize"
        @pagination="getAccountList"
      />
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="dialogMode === 'create' ? '新建账号' : '编辑账号'"
      width="500px"
      @close="handleCloseDialog"
    >
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="formData.username"
            placeholder="请输入用户名"
            :disabled="dialogMode === 'edit'"
          />
        </el-form-item>
        <el-form-item label="真实姓名" prop="realName">
          <el-input v-model="formData.realName" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="手机号" prop="mobile">
          <el-input v-model="formData.mobile" placeholder="请输入手机号" />
        </el-form-item>
        <template v-if="dialogMode === 'create'">
          <el-form-item label="密码" prop="password">
            <el-input
              v-model="formData.password"
              type="password"
              show-password
              placeholder="请输入密码"
            />
          </el-form-item>
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input
              v-model="formData.confirmPassword"
              type="password"
              show-password
              placeholder="请再次输入新密码"
            />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit"> 确定 </el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog
      v-model="resetPwdDialogVisible"
      title="重置密码"
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
          <el-button type="primary" @click="handleResetPwdSubmit"> 确定 </el-button>
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
    createAccountApi,
    listAccountsApi,
    updateAccountStatusApi,
    deleteAccountApi,
    updateAccountApi,
    resetPasswordApi,
    type CreateAccountPayload,
    type UpdateAccountPayload,
    type ResetPasswordPayload,
  } from '@/api/account';

  interface SysUser {
    ID: number;
    Username: string;
    RealName: string;
    Mobile: string;
    Status: 1 | 2;
    CreatedAt: string;
  }

  // --- 列表和分页 ---
  const loading = ref(false);
  const accountList = ref<SysUser[]>([]);
  const total = ref(0);
  const pagination = reactive({
    page: 1,
    pageSize: 20, // 默认每页20条
  });

  const getAccountList = async () => {
    loading.value = true;
    try {
      const res = await listAccountsApi({
        page: pagination.page,
        pageSize: pagination.pageSize,
      });
      accountList.value = res.list;
      total.value = res.total;
    } catch (error) {
      console.error(error);
      accountList.value = [];
      total.value = 0;
    } finally {
      loading.value = false;
    }
  };

  onMounted(() => {
    getAccountList();
  });

  // --- 新建/编辑对话框 ---
  const dialogVisible = ref(false);
  const dialogMode = ref<'create' | 'edit'>('create');
  const currentAccountId = ref<number | null>(null);
  const formRef = ref<FormInstance | null>(null);

  const getInitialFormData = (): CreateAccountPayload & {
    confirmPassword: '';
  } & UpdateAccountPayload => ({
    username: '',
    password: '',
    confirmPassword: '',
    realName: '',
    mobile: '',
  });

  const formData = reactive(getInitialFormData());

  const validatePass = (_rule: any, value: any, callback: any) => {
    if (value === '') {
      callback(new Error('请再次输入密码'));
    } else if (value !== formData.password) {
      callback(new Error('两次输入的密码不一致'));
    } else {
      callback();
    }
  };

  const formRules = reactive<FormRules>({
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    realName: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
    password: [
      { required: true, message: '请输入密码', trigger: 'blur' },
      { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
    ],
    confirmPassword: [{ required: true, validator: validatePass, trigger: 'blur' }],
  });

  const handleOpenDialog = () => {
    dialogMode.value = 'create';
    dialogVisible.value = true;
  };

  const handleEdit = (row: SysUser) => {
    dialogMode.value = 'edit';
    currentAccountId.value = row.ID;
    formData.username = row.Username;
    formData.realName = row.RealName;
    formData.mobile = row.Mobile;
    dialogVisible.value = true;
  };

  const handleCloseDialog = () => {
    formRef.value?.resetFields();
    Object.assign(formData, getInitialFormData());
    currentAccountId.value = null;
  };

  const handleSubmit = async () => {
    if (!formRef.value) return;
    await formRef.value.validate(async (valid) => {
      if (valid) {
        try {
          if (dialogMode.value === 'create') {
            const payload: CreateAccountPayload = {
              username: formData.username,
              realName: formData.realName,
              password: formData.password,
              mobile: formData.mobile,
            };
            await createAccountApi(payload);
            ElMessage.success('账号创建成功');
          } else {
            if (!currentAccountId.value) return;
            const payload: UpdateAccountPayload = {
              realName: formData.realName,
              mobile: formData.mobile,
            };
            await updateAccountApi(currentAccountId.value, payload);
            ElMessage.success('账号更新成功');
          }
          dialogVisible.value = false;
          await getAccountList();
        } catch (error) {
          console.error(error);
        }
      }
    });
  };

  // --- 重置密码对话框 ---
  const resetPwdDialogVisible = ref(false);
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

  const handleOpenResetPwdDialog = (row: SysUser) => {
    currentAccountId.value = row.ID;
    resetPwdDialogVisible.value = true;
  };

  const handleCloseResetPwdDialog = () => {
    resetPwdFormRef.value?.resetFields();
    currentAccountId.value = null;
  };

  const handleResetPwdSubmit = async () => {
    if (!resetPwdFormRef.value || !currentAccountId.value) return;
    await resetPwdFormRef.value.validate(async (valid) => {
      if (valid) {
        try {
          const payload: ResetPasswordPayload = { password: resetPwdFormData.password };
          await resetPasswordApi(currentAccountId.value!, payload);
          ElMessage.success('密码重置成功');
          resetPwdDialogVisible.value = false;
        } catch (error) {
          console.error(error);
        }
      }
    });
  };

  // --- 其他操作 ---
  const handleStatusChange = async (row: SysUser) => {
    const newStatus = row.Status;
    const oldStatus = newStatus === 1 ? 2 : 1;
    const actionText = newStatus === 1 ? '启用' : '禁用';

    try {
      await ElMessageBox.confirm(
        `确定要<strong style="color: #F56C6C">${actionText}</strong>用户 [${row.Username}] 吗？`,
        '提示',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
          dangerouslyUseHTMLString: true,
        }
      );

      await updateAccountStatusApi(row.ID, { status: newStatus });
      ElMessage.success(`${actionText}成功`);
    } catch (error) {
      // @ts-ignore
      row.Status = oldStatus;
      if (error !== 'cancel') {
        ElMessage.error('操作失败');
      }
    }
  };

  const handleDelete = async (row: SysUser) => {
    try {
      await ElMessageBox.confirm(
        `确定要永久<strong style="color: #F56C6C">删除</strong>用户 [${row.Username}] 吗？此操作不可恢复。`,
        '危险操作',
        {
          confirmButtonText: '确定删除',
          cancelButtonText: '取消',
          type: 'error',
          dangerouslyUseHTMLString: true,
        }
      );

      await deleteAccountApi(row.ID);
      ElMessage.success('删除成功');
      await getAccountList();
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error('删除失败');
      }
    }
  };
</script>

<style scoped>
  .app-container {
    padding: 20px;
    display: flex;
    flex-direction: column;
    height: calc(100vh - 90px); /* 减去顶部导航和内边距 */
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
    flex-grow: 1; /* 让表格占据多余空间 */
    overflow-y: auto;
  }

  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
    flex-shrink: 0; /* 防止分页组件被压缩 */
  }
</style>
