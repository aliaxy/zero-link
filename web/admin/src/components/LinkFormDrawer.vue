<template>
  <el-drawer
    v-model="visible"
    :title="isEdit ? 'Edit Link' : 'New Link'"
    direction="rtl"
    size="480px"
    @closed="resetForm"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-position="top"
      class="link-form"
    >
      <el-form-item label="Destination URL" prop="origin_url">
        <el-input
          v-model="form.origin_url"
          placeholder="https://example.com/page"
          :disabled="isEdit"
        />
      </el-form-item>

      <el-form-item label="Custom Code" prop="code">
        <el-input
          v-model="form.code"
          placeholder="e.g. my-link (optional)"
          :disabled="isEdit"
        >
          <template #prefix>
            <span class="code-prefix">·/</span>
          </template>
        </el-input>
        <div class="field-hint">3–32 chars, letters, digits, - and _. Leave blank to auto-generate.</div>
      </el-form-item>

      <el-form-item label="Title" prop="title">
        <el-input v-model="form.title" placeholder="Optional label for this link" />
      </el-form-item>

      <el-form-item label="Description" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="2"
          placeholder="Optional note"
        />
      </el-form-item>

      <el-form-item label="Expiry Date" prop="expire_at">
        <el-date-picker
          v-model="form.expire_at"
          type="datetime"
          placeholder="No expiry"
          format="YYYY-MM-DD HH:mm"
          :disabled-date="(d) => d < new Date()"
          style="width: 100%"
        />
      </el-form-item>

      <el-form-item v-if="isEdit" label="Status" prop="status">
        <el-select v-model="form.status" style="width: 100%">
          <el-option label="Active" :value="1" />
          <el-option label="Disabled" :value="0" />
        </el-select>
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="drawer-footer">
        <el-button @click="visible = false">Cancel</el-button>
        <el-button type="primary" :loading="loading" @click="submit">
          {{ isEdit ? 'Save Changes' : 'Create Link' }}
        </el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { createLink, updateLink } from '@/api/links'
import dayjs from 'dayjs'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  link: { type: Object, default: null },
})

const emit = defineEmits(['update:modelValue', 'saved'])

const visible = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v),
})

const isEdit = computed(() => !!props.link)
const formRef = ref(null)
const loading = ref(false)

const defaultForm = () => ({
  origin_url: '',
  code: '',
  title: '',
  description: '',
  expire_at: null,
  status: 1,
})

const form = ref(defaultForm())

const RESERVED = new Set(['admin', 'api', 'healthz', 'readyz', 'metrics', 'static'])
const CODE_RE = /^[a-zA-Z0-9_-]{3,32}$/

const rules = {
  origin_url: [
    { required: true, message: 'URL is required', trigger: 'blur' },
    {
      validator: (_, v, cb) => {
        try {
          const u = new URL(v)
          if (!['http:', 'https:'].includes(u.protocol)) cb(new Error('Must be http or https'))
          else cb()
        } catch {
          cb(new Error('Enter a valid URL'))
        }
      },
      trigger: 'blur',
    },
  ],
  code: [
    {
      validator: (_, v, cb) => {
        if (!v) { cb(); return }
        if (!CODE_RE.test(v)) { cb(new Error('3–32 chars: letters, digits, - and _')); return }
        if (RESERVED.has(v.toLowerCase())) { cb(new Error('This code is reserved')); return }
        cb()
      },
      trigger: 'blur',
    },
  ],
}

watch(
  () => props.link,
  (link) => {
    if (link) {
      form.value = {
        origin_url: link.origin_url,
        code: link.code,
        title: link.title ?? '',
        description: link.description ?? '',
        expire_at: link.expire_at ? new Date(link.expire_at) : null,
        status: link.status,
      }
    } else {
      form.value = defaultForm()
    }
  },
  { immediate: true },
)

function resetForm() {
  form.value = defaultForm()
  formRef.value?.clearValidate()
}

async function submit() {
  await formRef.value.validate()
  loading.value = true
  try {
    const payload = {
      origin_url: form.value.origin_url,
      title: form.value.title || undefined,
      description: form.value.description || undefined,
      expire_at: form.value.expire_at
        ? dayjs(form.value.expire_at).toISOString()
        : undefined,
    }
    if (!isEdit.value && form.value.code) {
      payload.code = form.value.code
    }
    if (isEdit.value) {
      payload.status = form.value.status
      await updateLink(props.link.id, payload)
    } else {
      await createLink(payload)
    }
    ElMessage.success(isEdit.value ? 'Link updated' : 'Link created')
    visible.value = false
    emit('saved')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.link-form {
  padding: 4px 0;
}

.code-prefix {
  font-family: var(--font-mono);
  font-size: 13px;
  color: var(--color-label-tertiary);
  margin-right: 2px;
}

.field-hint {
  font-size: 12px;
  color: var(--color-label-tertiary);
  margin-top: 5px;
  line-height: 1.4;
}

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding-top: 4px;
}
</style>
