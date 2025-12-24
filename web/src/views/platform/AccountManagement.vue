<template>
  <div class="app-container">
    <div class="toolbar">
      <el-button type="primary" :icon="Plus" @click="handleOpenDialog">
        新建账号
      </el-button>
    </div>

    <!-- 账号列表 -->
    <el-table v-loading="loading" :data="accountList" border stripe style="width: 100%; margin-top: 20px;">
      <el-table-column prop="ID" label="ID" width="80" />
      <el-table-column prop="Username" label="用户名" />
      <el-table-column prop="RealName" label="真实姓名" />
      <el-table-column prop="Mobile" label="手机号" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-switch 
            v-model="row.Status"
            :active-value="1"
            :inactive-value="2"
            @change="handleStatusChange(row)"
          />
        </template>
      </el-table-column>
      <el-table-column prop="CreatedAt" label="创建时间" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button type="primary" link>编辑</el-button>
          <el-button type="primary" link>重置密码</el-button>
          <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination-container">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 新建/编辑账号对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      title="新建平台账号"
      width="500px"
      @close="handleCloseDialog"
    >
      <el-form 
        ref="formRef"
        :model="formData" 
        :rules="formRules" 
        label-width="100px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="真实姓名" prop="realName">
          <el-input v-model="formData.realName" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="手机号" prop="mobile">
          <el-input v-model="formData.mobile" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="formData.password" type="password" show-password placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="formData.confirmPassword" type="password" show-password placeholder="请再次输入密码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">
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
import { Plus } from '@element-plus/icons-vue';
import type { FormInstance, FormRules } from 'element-plus';
import { createAccountApi, listAccountsApi, updateAccountStatusApi, deleteAccountApi, type CreateAccountPayload } from '@/api/account';

// --- 列表和分页 ---
const loading = ref(false);
const accountList = ref<any[]>([]); // 明确类型
const total = ref(0);
const pagination = reactive({
  page: 1,
  pageSize: 10,
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

const handleSizeChange = (size: number) => {
  pagination.pageSize = size;
  getAccountList();
};

const handleCurrentChange = (page: number) => {
  pagination.page = page;
  getAccountList();
};


// --- 对话框和表单 ---
const dialogVisible = ref(false);
const formRef = ref<FormInstance | null>(null);

const getInitialFormData = (): CreateAccountPayload & { confirmPassword: '' } => ({
  username: '',
  password: '',
  confirmPassword: '',
  realName: '',
  mobile: '',
});

const formData = reactive(getInitialFormData());

// --- 表单校验规则 ---
const validatePass = (_rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请再次输入密码'));
  } else if (value !== formData.password) {
    callback(new Error("两次输入的密码不一致"));
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
  confirmPassword: [
    { required: true, validator: validatePass, trigger: 'blur' }
  ],
});


// --- 事件处理 ---
const handleOpenDialog = () => {
  dialogVisible.value = true;
};

const handleCloseDialog = () => {
  formRef.value?.resetFields();
  Object.assign(formData, getInitialFormData());
};

const handleSubmit = async () => {
  if (!formRef.value) return;
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const payload: CreateAccountPayload = {
          username: formData.username,
          realName: formData.realName,
          password: formData.password,
          mobile: formData.mobile,
        };
        await createAccountApi(payload);
        ElMessage.success('账号创建成功');
        dialogVisible.value = false;
        await getAccountList(); // 刷新账号列表
      } catch (error) {
        console.error(error);
      }
    }
  });
};

const handleStatusChange = async (row: any) => {
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
    // 状态已由 v-model 更新，无需手动刷新列表
  } catch (error) {
    // 如果是用户取消(cancel)或API调用失败
    row.Status = oldStatus; // 恢复开关到原始状态
    if (error !== 'cancel') {
      ElMessage.error('操作失败');
    }
  }
};

const handleDelete = async (row: any) => {
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
    await getAccountList(); // 刷新列表
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
}
.toolbar {
  margin-bottom: 20px;
}
.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
