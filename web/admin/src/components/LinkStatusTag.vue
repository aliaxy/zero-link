<template>
  <el-tag :type="tagType" size="small" effect="light">{{ label }}</el-tag>
</template>

<script setup>
import { computed } from 'vue'
import dayjs from 'dayjs'

const props = defineProps({
  status: { type: Number, required: true },
  expireAt: { type: String, default: null },
})

const isExpired = computed(() =>
  props.expireAt && dayjs(props.expireAt).isBefore(dayjs()),
)

const tagType = computed(() => {
  if (isExpired.value) return 'info'
  if (props.status === 2) return 'danger'
  return 'success'
})

const label = computed(() => {
  if (isExpired.value) return 'Expired'
  if (props.status === 2) return 'Disabled'
  return 'Active'
})
</script>
