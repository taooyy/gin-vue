<template>
  <el-pagination
    v-if="!hidden"
    v-model:current-page="currentPage"
    v-model:page-size="pageSize"
    :background="background"
    :layout="layout"
    :page-sizes="pageSizes"
    :total="total"
    @size-change="handleSizeChange"
    @current-change="handleCurrentChange"
  />
</template>

<script setup lang="ts">
  import { computed } from 'vue';

  const props = defineProps({
    total: {
      type: Number,
      required: true,
    },
    page: {
      type: Number,
      default: 1,
    },
    limit: {
      type: Number,
      default: 20,
    },
    pageSizes: {
      type: Array as () => number[],
      default: () => [10, 20, 30, 50],
    },
    layout: {
      type: String,
      default: 'total, sizes, prev, pager, next, jumper',
    },
    background: {
      type: Boolean,
      default: true,
    },
    hidden: {
      type: Boolean,
      default: false,
    },
  });

  const emit = defineEmits(['update:page', 'update:limit', 'pagination']);

  const currentPage = computed({
    get() {
      return props.page;
    },
    set(val) {
      emit('update:page', val);
    },
  });

  const pageSize = computed({
    get() {
      return props.limit;
    },
    set(val) {
      emit('update:limit', val);
    },
  });

  const handleSizeChange = (val: number) => {
    emit('pagination', { page: currentPage.value, limit: val });
  };

  const handleCurrentChange = (val: number) => {
    emit('pagination', { page: val, limit: pageSize.value });
  };
</script>
