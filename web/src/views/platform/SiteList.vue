<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>站点展示</span>
        </div>
      </template>

      <div v-loading="loading" class="card-grid-container">
        <div v-if="schoolList.length > 0" class="card-grid">
          <el-row :gutter="20">
            <el-col v-for="school in schoolList" :key="school.ID" :xs="24" :sm="12" :md="8" :lg="6">
              <el-card class="school-card" shadow="hover">
                <template #header>
                  <div class="school-card-header">
                    <span>{{ school.Name }}</span>
                    <el-tag :type="school.IsEnabled ? 'success' : 'info'" size="small">
                      {{ school.IsEnabled ? '运营中' : '已禁用' }}
                    </el-tag>
                  </div>
                </template>
                <div class="school-card-body">
                  <p>
                    <el-icon><User /></el-icon> 负责人: {{ school.ContactName }}
                  </p>
                  <p>
                    <el-icon><Phone /></el-icon> 联系电话: {{ school.ContactPhone }}
                  </p>
                  <p>
                    <el-icon><Location /></el-icon> 地址: {{ school.Address }}
                  </p>
                </div>
                <div class="school-card-footer">
                  <el-button type="primary" link @click="goToManage()">详情与管理</el-button>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </div>
        <el-empty v-else description="暂无站点数据" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue';
  import { useRouter } from 'vue-router';
  import { listSchoolsApi, type School } from '@/api/school';
  import { User, Phone, Location } from '@element-plus/icons-vue';

  const loading = ref(false);
  const schoolList = ref<School[]>([]);
  const router = useRouter();

  const getList = async () => {
    loading.value = true;
    try {
      const res = await listSchoolsApi({ page: 1, pageSize: 999 });
      schoolList.value = res.list;
    } catch (error) {
      console.error(error);
    } finally {
      loading.value = false;
    }
  };

  const goToManage = () => {
    router.push({ name: 'SiteManage' });
  };

  onMounted(() => {
    getList();
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
    flex-grow: 1;
    overflow-y: auto;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .school-card {
    margin-bottom: 20px;
    display: flex;
    flex-direction: column;
    height: 280px;
  }

  .school-card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: bold;
  }

  .school-card-body {
    flex-grow: 1;
    font-size: 14px;
    color: #606266;
  }

  .school-card-body p {
    margin: 0 0 12px 0;
    display: flex;
    align-items: center;
  }

  .school-card-body .el-icon {
    margin-right: 8px;
    font-size: 16px;
  }

  .school-card-footer {
    border-top: 1px solid #e4e7ed;
    padding-top: 10px;
    margin-top: 10px;
    text-align: right;
    flex-shrink: 0;
  }
</style>
