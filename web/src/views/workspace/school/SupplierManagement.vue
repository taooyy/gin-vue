<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>供应商管理</span>
          <el-button type="primary" :icon="Plus" @click="handleOpenDialog()"> 新建供应商 </el-button>
        </div>
      </template>

      <el-table v-loading="loading" :data="supplierList" border stripe style="width: 100%">
        <el-table-column prop="ID" label="ID" width="80" align="center" />
        <el-table-column prop="Name" label="供应商名称" min-width="150" />
        <el-table-column prop="ContactName" label="负责人" min-width="120" />
        <el-table-column prop="ContactPhone" label="联系电话" min-width="150" />
        <el-table-column prop="Address" label="地址" min-width="200" show-overflow-tooltip />
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.IsEnabled ? 'success' : 'info'">{{
              row.IsEnabled ? '启用' : '禁用'
            }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180" align="center">
          <template #default="{ row }">
            {{ new Date(row.CreatedAt).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleOpenDialog(row.ID)">
              编辑
            </el-button>
            <el-button
              :type="row.IsEnabled ? 'danger' : 'success'"
              link
              size="small"
              @click="handleUpdateStatus(row)"
            >
              {{ row.IsEnabled ? '禁用' : '启用' }}
            </el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" @close="handleCloseDialog">
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px">
        <el-form-item label="供应商名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入供应商公司名称" />
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
        <el-divider>管理员信息</el-divider>
        <el-form-item label="管理员账号" prop="username">
          <el-input
            v-model="formData.username"
            placeholder="用于登录供应商后台的管理员账号"
            :disabled="!!editingID"
          />
        </el-form-item>
        <el-form-item label="管理员姓名" prop="realName">
          <el-input v-model="formData.realName" placeholder="请输入管理员真实姓名" />
        </el-form-item>
        <template v-if="!editingID">
          <el-form-item label="初始密码" prop="password">
            <el-input
              v-model="formData.password"
              type="password"
              show-password
              placeholder="请输入初始密码"
            />
          </el-form-item>
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input
              v-model="formData.confirmPassword"
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue';
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';
import Pagination from '@/components/Pagination/index.vue';
import {
  listSuppliersApi,
  createSupplierApi,
  getSupplierApi,
  updateSupplierApi,
  updateSupplierStatusApi,
  type Supplier,
  type CreateSupplierPayload,
} from '@/api/supplier';

// --- 状态引用 ---
const isSubmitting = ref(false);

// --- 列表和分页 ---
const loading = ref(false);
const supplierList = ref<Supplier[]>([]);
const total = ref(0);
const pagination = reactive({
  page: 1,
  pageSize: 20,
});

/**
 * @description 获取供应商列表数据
 */
const getList = async () => {
  loading.value = true;
  try {
    const res = await listSuppliersApi(pagination);
    supplierList.value = res.items;
    total.value = res.total;
  } catch (error) {
    console.error(error);
  } finally {
    loading.value = false;
  }
};

// --- 生命周期钩子 ---
onMounted(() => {
  getList();
});

const editingID = ref<number | null>(null);

// --- 新建/编辑对话框 ---
const dialogVisible = ref(false);
const dialogTitle = computed(() => (editingID.value ? '编辑供应商' : '新建供应商'));
const formRef = ref<FormInstance | null>(null);

const getInitialFormData = (): CreateSupplierPayload & { confirmPassword: '' } => ({
  name: '',
  contactName: '',
  contactPhone: '',
  address: '',
  username: '',
  password: '',
  realName: '',
  confirmPassword: '',
});

const formData = reactive(getInitialFormData());

// --- 表单校验规则 ---
const validatePass = (_rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请再次输入密码'));
  } else if (value !== formData.password) {
    callback(new Error('两次输入的密码不一致'));
  } else {
    callback();
  }
};

const formRules = computed<FormRules>(() => {
  const rules: FormRules = {
    name: [{ required: true, message: '请输入供应商名称', trigger: 'blur' }],
    contactName: [{ required: true, message: '请输入负责人姓名', trigger: 'blur' }],
    username: [{ required: true, message: '请输入管理员账号', trigger: 'blur' }],
    realName: [{ required: true, message: '请输入管理员姓名', trigger: 'blur' }],
  };

  if (!editingID.value) {
    rules.password = [
      { required: true, message: '请输入初始密码', trigger: 'blur' },
      { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
    ];
    rules.confirmPassword = [{ required: true, validator: validatePass, trigger: 'blur' }];
  }

  return rules;
});

/**
 * @description 打开新建/编辑对话框
 */
const handleOpenDialog = async (id?: number) => {
  if (id) {
    editingID.value = id;
    try {
      const supplierDetails = await getSupplierApi(id);
      Object.assign(formData, {
        name: supplierDetails.Name,
        contactName: supplierDetails.ContactName,
        contactPhone: supplierDetails.ContactPhone,
        address: supplierDetails.Address,
        username: supplierDetails.adminUser?.Username || '',
        realName: supplierDetails.adminUser?.RealName || '',
        // Clear password fields for editing
        password: '',
        confirmPassword: '',
      });
    } catch (error) {
      console.error('Failed to fetch supplier details', error);
      ElMessage.error('获取供应商信息失败');
      return;
    }
  } else {
    editingID.value = null;
    // Reset form to initial state for creation
    Object.assign(formData, getInitialFormData());
    // Ensure validation state is also fresh
    formRef.value?.clearValidate();
  }
  dialogVisible.value = true;
};

/**
 * @description 关闭对话框时的回调，重置表单和编辑状态
 */
const handleCloseDialog = () => {
  editingID.value = null;
  formRef.value?.resetFields();
  Object.assign(formData, getInitialFormData());
  dialogVisible.value = false;
};

/**
 * @description 处理新建/更新表单的提交
 */
const handleSubmit = async () => {
  if (!formRef.value) return;
  await formRef.value.validate(async (valid) => {
    if (valid) {
      isSubmitting.value = true;
      try {
        if (editingID.value) {
          // Update logic - only send fields that are editable
          const updatePayload = {
            name: formData.name,
            contactName: formData.contactName,
            contactPhone: formData.contactPhone,
            address: formData.address,
            realName: formData.realName,
          };
          await updateSupplierApi(editingID.value, updatePayload);
          ElMessage.success('供应商更新成功');
        } else {
          // Create logic
          await createSupplierApi(formData);
          ElMessage.success('供应商创建成功');
        }
        handleCloseDialog();
        await getList(); // 刷新列表
      } catch (error) {
        console.error(error);
        // Can add error message display here
      } finally {
        isSubmitting.value = false;
      }
    }
  });
};

/**
 * @description 更新供应商的启用/禁用状态
 */
const handleUpdateStatus = async (row: Supplier) => {
  const actionText = row.IsEnabled ? '禁用' : '启用';
  try {
    await ElMessageBox.confirm(`确定要${actionText}供应商 "${row.Name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    });

    await updateSupplierStatusApi(row.ID, !row.IsEnabled);
    ElMessage.success(`${actionText}成功`);
    await getList(); // 刷新列表
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error);
    }
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
