<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>操作日志</span>
        </div>
      </template>

      <div class="filter-container">
        <el-select
          v-model="selectedOrgId"
          placeholder="请选择站点"
          clearable
          @change="handleFilter"
        >
          <el-option label="平台日志" :value="1" />
          <el-option
            v-for="item in schoolList"
            :key="item.ID"
            :label="item.Name"
            :value="item.ID"
          />
        </el-select>
      </div>

      <el-table v-loading="loading" :data="logList" border stripe style="width: 100%">
        <el-table-column prop="ID" label="ID" width="80" align="center" />
        <el-table-column prop="Username" label="操作人" min-width="120" />
        <el-table-column prop="Module" label="操作模块" min-width="200" show-overflow-tooltip />
        <el-table-column prop="Action" label="操作类型" width="100" align="center" />
        <el-table-column prop="CreatedAt" label="操作时间" min-width="180" />
        <el-table-column label="操作内容" min-width="300">
          <template #default="{ row }">
            <pre class="params-pre">{{ formatParams(row.Params) }}</pre>
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
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted } from 'vue';
  import Pagination from '@/components/Pagination/index.vue';
  import { listLogsApi, type OpLog, type ListLogsParams } from '@/api/log';
  import { listSchoolsApi, type School } from '@/api/school';

  // --- 响应式状态定义 ---
  const loading = ref(false); // 表格加载状态
  const logList = ref<OpLog[]>([]); // 日志列表数据
  const schoolList = ref<School[]>([]); // 用于下拉选择的学校列表
  const total = ref(0); // 日志总条数
  const selectedOrgId = ref<number | undefined>(1); // 选中的组织ID，默认为1（平台）

  // 分页数据
  const pagination = reactive({
    page: 1,
    pageSize: 20,
  });

  // --- API 调用封装 ---

  /**
   * @description 获取日志列表。
   * 根据 selectedOrgId 筛选特定站点的日志，或获取平台日志。
   */
  const getList = async () => {
    loading.value = true;
    const params: ListLogsParams = {
      page: pagination.page,
      pageSize: pagination.pageSize,
    };
    // 如果 selectedOrgId 有值，则作为查询参数
    if (selectedOrgId.value) {
      params.orgId = selectedOrgId.value;
    }

    try {
      const res = await listLogsApi(params);
      logList.value = res.list;
      total.value = res.total;
    } catch (error) {
      console.error('获取日志列表失败:', error);
    } finally {
      loading.value = false;
    }
  };

  /**
   * @description 获取所有学校的列表，用于填充筛选下拉框。
   */
  const getSchoolList = async () => {
    try {
      // 设置一个很大的 pageSize 以便一次性获取所有学校
      const res = await listSchoolsApi({ page: 1, pageSize: 999 });
      schoolList.value = res.list;
    } catch (error) {
      console.error('获取学校列表失败:', error);
    }
  };

  // --- 事件处理 ---

  /**
   * @description 当用户在下拉框中选择不同站点时触发，重置页码并重新获取数据。
   */
  const handleFilter = () => {
    pagination.page = 1;
    getList();
  };

  /**
   * @description 格式化参数字符串为易读的 JSON 格式。
   * @param params - 后端返回的原始参数字符串
   */
  const formatParams = (params: string) => {
    if (!params) return '无参数';
    try {
      const obj = JSON.parse(params);
      return JSON.stringify(obj, null, 2); // 格式化为带缩进的JSON字符串
    } catch (e) {
      return params; // 如果解析失败，则直接返回原始字符串
    }
  };

  // --- 生命周期钩子 ---

  // 组件挂载后，获取初始的日志列表和学校列表
  onMounted(() => {
    getList();
    getSchoolList();
  });
</script>

<style scoped>
  .app-container {
    height: 100%;
  }
  .el-card {
    height: 100%;
    display: flex;
    flex-direction: column;
  }
  :deep(.el-card__body) {
    display: flex;
    flex-direction: column;
    flex-grow: 1;
    overflow: hidden;
  }
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .filter-container {
    margin-bottom: 20px;
  }
  .el-table {
    flex-grow: 1;
    overflow-y: auto;
  }
  .pagination-container {
    margin-top: 20px;
    flex-shrink: 0;
    display: flex;
    justify-content: flex-end;
  }
  .params-pre {
    white-space: pre-wrap;
    word-break: break-all;
    margin: 0;
    font-family: Menlo, Monaco, 'Courier New', monospace;
    font-size: 12px;
  }
</style>
