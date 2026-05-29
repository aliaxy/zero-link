<template>
  <div>
    <div class="back-row">
      <el-button text :icon="ArrowLeft" @click="router.push('/links')">Links</el-button>
    </div>

    <div v-if="link" class="detail-layout">
      <!-- Info card -->
      <el-card class="info-card">
        <div class="info-header">
          <span class="info-code mono">{{ link.code }}</span>
          <LinkStatusTag :status="link.status" :expire-at="link.expire_at" />
        </div>

        <div class="short-url-row">
          <span class="short-url mono">{{ shortBase }}/{{ link.code }}</span>
          <el-tooltip content="Copy short URL">
            <el-button text :icon="CopyDocument" @click="copyShortUrl" />
          </el-tooltip>
          <el-tooltip content="View stats">
            <el-button text :icon="DataLine" @click="router.push(`/links/${link.id}/stats`)" />
          </el-tooltip>
        </div>

        <div class="info-rows">
          <div class="info-row">
            <span class="info-label">Destination</span>
            <a :href="link.origin_url" target="_blank" rel="noopener" class="info-value link-value">
              {{ link.origin_url }}
            </a>
          </div>
          <div v-if="link.title" class="info-row">
            <span class="info-label">Title</span>
            <span class="info-value">{{ link.title }}</span>
          </div>
          <div v-if="link.description" class="info-row">
            <span class="info-label">Note</span>
            <span class="info-value">{{ link.description }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">Expires</span>
            <span class="info-value">{{ link.expire_at ? formatDate(link.expire_at) : 'Never' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">Created</span>
            <span class="info-value">{{ formatDate(link.created_at) }}</span>
          </div>
        </div>
      </el-card>

      <!-- Edit card -->
      <el-card class="edit-card">
        <template #header>
          <span class="edit-title">Edit Link</span>
        </template>

        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-position="top"
        >
          <el-form-item label="Destination URL" prop="origin_url">
            <el-input v-model="form.origin_url" placeholder="https://example.com" />
          </el-form-item>
          <el-form-item label="Title">
            <el-input v-model="form.title" placeholder="Optional" />
          </el-form-item>
          <el-form-item label="Description">
            <el-input v-model="form.description" type="textarea" :rows="2" placeholder="Optional" />
          </el-form-item>
          <el-form-item label="Expiry Date">
            <el-date-picker
              v-model="form.expire_at"
              type="datetime"
              placeholder="No expiry"
              format="YYYY-MM-DD HH:mm"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="Status">
            <el-select v-model="form.status" style="width: 100%">
              <el-option label="Active" :value="1" />
              <el-option label="Disabled" :value="2" />
            </el-select>
          </el-form-item>

          <el-button type="primary" :loading="saving" @click="saveChanges">Save Changes</el-button>
        </el-form>
      </el-card>
    </div>

    <div v-else-if="notFound" class="empty-state">
      <p>Link not found.</p>
      <el-button @click="router.push('/links')">Back to Links</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, CopyDocument, DataLine } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getLink, updateLink } from '@/api/links'
import LinkStatusTag from '@/components/LinkStatusTag.vue'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()

const shortBase = import.meta.env.VITE_SHORT_LINK_BASE ?? window.location.origin

const link = ref(null)
const notFound = ref(false)
const saving = ref(false)
const formRef = ref(null)

const form = ref({
  origin_url: '',
  title: '',
  description: '',
  expire_at: null,
  status: 1,
})

const rules = {
  origin_url: [{ required: true, message: 'URL is required', trigger: 'blur' }],
}

onMounted(async () => {
  try {
    link.value = await getLink(route.params.id)
    form.value = {
      origin_url: link.value.origin_url,
      title: link.value.title ?? '',
      description: link.value.description ?? '',
      expire_at: link.value.expire_at ? new Date(link.value.expire_at) : null,
      status: link.value.status,
    }
  } catch {
    notFound.value = true
  }
})

async function saveChanges() {
  await formRef.value.validate()
  saving.value = true
  try {
    link.value = await updateLink(route.params.id, {
      origin_url: form.value.origin_url,
      title: form.value.title || undefined,
      description: form.value.description || undefined,
      expire_at: form.value.expire_at ? dayjs(form.value.expire_at).toISOString() : undefined,
      status: form.value.status,
    })
    ElMessage.success('Link updated')
  } finally {
    saving.value = false
  }
}

async function copyShortUrl() {
  await navigator.clipboard.writeText(`${shortBase}/${link.value.code}`)
  ElMessage.success('Copied!')
}

function formatDate(iso) {
  return dayjs(iso).format('MMM D, YYYY HH:mm')
}
</script>

<style scoped>
.back-row {
  margin-bottom: 20px;
}

.detail-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  align-items: start;
}

.info-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.info-code {
  font-size: 22px;
  font-weight: 600;
  color: var(--color-label);
  letter-spacing: -0.02em;
}

.short-url-row {
  display: flex;
  align-items: center;
  gap: 6px;
  background: var(--color-bg);
  border-radius: var(--radius-md);
  padding: 10px 14px;
  margin-bottom: 20px;
}

.short-url {
  flex: 1;
  font-size: 13px;
  color: var(--color-accent);
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.info-rows {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.info-row {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.info-label {
  font-size: 11px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--color-label-tertiary);
}

.info-value {
  font-size: 14px;
  color: var(--color-label);
  word-break: break-all;
}

.link-value {
  color: var(--color-accent);
  text-decoration: none;
}

.link-value:hover {
  text-decoration: underline;
}

.edit-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-label);
}

.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: var(--color-label-secondary);
}

@media (max-width: 900px) {
  .detail-layout {
    grid-template-columns: 1fr;
  }
}
</style>
