<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Check, ArrowLeft, Picture } from '@element-plus/icons-vue'
import { MdEditor } from 'md-editor-v3'
import type { ExposeParam } from 'md-editor-v3'
import { getPost, createPost, updatePost, deletePost } from '../api/posts'
import { uploadImage, listMedia } from '../api/upload'

const route = useRoute()
const router = useRouter()

const currentId = ref<number | null>(route.params.id ? Number(route.params.id) : null)
const editorRef = ref<ExposeParam>()

const title = ref('')
const summary = ref('')
const content = ref('')
const coverImage = ref('')
const coverInputRef = ref<HTMLInputElement>()
const saving = ref(false)
const autoCreated = ref(false)
const showMedia = ref(false)
const mediaFiles = ref<string[]>([])
const uploading = ref(false)
let dragCounter = 0
const draggingOver = ref(false)

onMounted(async () => {
  if (currentId.value) {
    const { data } = await getPost(currentId.value)
    title.value = data.title
    summary.value = data.summary
    content.value = data.content
    coverImage.value = data.cover_image || ''
    loadMedia()
  }
})

function payload() {
  return {
    title: title.value.trim(),
    summary: summary.value.trim(),
    content: content.value,
    cover_image: coverImage.value,
  }
}

async function onCoverChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  if (!await ensureId()) return
  uploading.value = true
  try {
    const { data } = await uploadImage(currentId.value!, file)
    coverImage.value = data.url
    await loadMedia()
  } finally {
    uploading.value = false
  }
}

async function ensureId(): Promise<boolean> {
  if (currentId.value) return true
  if (!title.value.trim()) {
    ElMessage.warning('Enter a title before uploading files')
    return false
  }
  const { data } = await createPost(payload())
  currentId.value = data.id
  autoCreated.value = true
  router.replace(`/posts/${data.id}/edit`)
  return true
}

async function loadMedia() {
  if (!currentId.value) return
  const { data } = await listMedia(currentId.value)
  mediaFiles.value = data.files
}

function toggleMedia() {
  showMedia.value = !showMedia.value
  if (showMedia.value && currentId.value) loadMedia()
}

async function handleFiles(files: File[]) {
  if (!await ensureId()) return
  uploading.value = true
  try {
    for (const f of files) {
      const { data } = await uploadImage(currentId.value!, f)
      insertUrl(data.url)
    }
    await loadMedia()
  } finally {
    uploading.value = false
  }
}

function insertUrl(url: string) {
  const name = url.split('/').pop()!
  const isImg = /\.(jpe?g|png|gif|webp|svg|bmp|ico)$/i.test(name)
  editorRef.value?.insert(() => ({
    targetValue: isImg ? `![${name}](${url})` : `[${name}](${url})`,
    select: false,
    deviationStart: 0,
    deviationEnd: 0,
  }))
}

// Used by md-editor-v3 toolbar image button only
async function onUploadImg(files: File[], callback: (urls: string[]) => void) {
  if (!await ensureId()) { callback([]); return }
  uploading.value = true
  try {
    const urls = await Promise.all(
      files.map(f => uploadImage(currentId.value!, f).then(r => r.data.url))
    )
    callback(urls)
    await loadMedia()
  } finally {
    uploading.value = false
  }
}

// Intercept ALL file pastes before md-editor-v3 sees them (capture phase)
async function onPasteCapture(e: ClipboardEvent) {
  if (!e.clipboardData?.files.length) return
  e.preventDefault()
  e.stopPropagation()
  await handleFiles(Array.from(e.clipboardData.files))
}

function onDragEnter() {
  dragCounter++
  draggingOver.value = true
}

function onDragLeave() {
  dragCounter--
  if (dragCounter <= 0) {
    dragCounter = 0
    draggingOver.value = false
  }
}

async function onDrop(e: DragEvent) {
  dragCounter = 0
  draggingOver.value = false
  const files = Array.from(e.dataTransfer?.files ?? [])
  if (files.length) await handleFiles(files)
}

async function onFileInput(e: Event) {
  const input = e.target as HTMLInputElement
  const files = Array.from(input.files ?? [])
  input.value = ''
  await handleFiles(files)
}

function isImage(url: string) {
  return /\.(jpe?g|png|gif|webp|svg|bmp|ico)$/i.test(url)
}

async function save() {
  if (!title.value.trim()) {
    ElMessage.warning('Title is required')
    return
  }
  saving.value = true
  try {
    if (currentId.value) {
      await updatePost(currentId.value, payload())
      ElMessage.success('Post saved')
      router.push(`/posts/${currentId.value}`)
    } else {
      const { data } = await createPost(payload())
      ElMessage.success('Post created')
      router.push(`/posts/${data.id}`)
    }
  } finally {
    saving.value = false
  }
}

async function cancel() {
  if (currentId.value && autoCreated.value) {
    await deletePost(currentId.value).catch(() => {})
  }
  router.back()
}
</script>

<template>
  <div>
    <div class="page-header">
      <h2>{{ currentId ? 'Edit Post' : 'New Post' }}</h2>
      <div style="display: flex; gap: 8px; align-items: center">
        <el-button :icon="Picture" :type="showMedia ? 'primary' : ''" @click="toggleMedia">
          Media<span v-if="mediaFiles.length"> ({{ mediaFiles.length }})</span>
        </el-button>
        <el-button :icon="ArrowLeft" @click="cancel">Cancel</el-button>
        <el-button type="primary" :icon="Check" :loading="saving" @click="save">Save</el-button>
      </div>
    </div>

    <el-card shadow="never" style="border-radius: 10px">
      <el-form label-position="top" class="editor-form">
        <el-form-item label="Cover Image">
          <div
            class="cover-zone"
            :class="{ 'has-cover': coverImage }"
            @click="coverInputRef?.click()"
          >
            <img v-if="coverImage" :src="coverImage" class="cover-preview" alt="" />
            <span v-else class="cover-hint">Click to add cover image</span>
            <button v-if="coverImage" class="cover-remove" @click.stop="coverImage = ''" title="Remove">×</button>
          </div>
          <input ref="coverInputRef" type="file" accept="image/*" @change="onCoverChange" style="display: none" />
        </el-form-item>

        <el-form-item label="Title" required>
          <el-input v-model="title" placeholder="Post title" size="large" />
        </el-form-item>

        <el-form-item label="Summary">
          <el-input v-model="summary" placeholder="Brief summary (optional)" />
        </el-form-item>

        <el-form-item label="Content">
          <div
            class="editor-drop-zone"
            :class="{ 'is-dragging': draggingOver }"
            @paste.capture="onPasteCapture"
            @dragenter.prevent="onDragEnter"
            @dragover.prevent
            @dragleave="onDragLeave"
            @drop.prevent="onDrop"
          >
            <MdEditor
              ref="editorRef"
              v-model="content"
              language="en-US"
              :onUploadImg="onUploadImg"
              style="width: 100%; height: 520px"
            />
            <div v-if="draggingOver || uploading" class="drop-mask">
              {{ uploading ? 'Uploading…' : 'Drop files to upload' }}
            </div>
          </div>
        </el-form-item>

        <!-- Media panel -->
        <div v-if="showMedia" class="media-panel">
          <div class="media-panel-header">
            <span v-if="currentId">
              {{ mediaFiles.length }} file{{ mediaFiles.length !== 1 ? 's' : '' }} in media/{{ currentId }}/
            </span>
            <span v-else style="color: var(--el-text-color-secondary)">
              Save the post first to manage media
            </span>
            <label v-if="currentId" class="media-upload-btn">
              + Upload files
              <input type="file" multiple @change="onFileInput" style="display: none" />
            </label>
          </div>
          <div v-if="currentId && !mediaFiles.length" class="media-empty">
            No files yet — paste, drag, or click Upload files above
          </div>
          <div v-if="currentId && mediaFiles.length" class="media-grid">
            <div
              v-for="url in mediaFiles"
              :key="url"
              class="media-item"
              :title="`Click to insert: ${url.split('/').pop()}`"
              @click="insertUrl(url)"
            >
              <img v-if="isImage(url)" :src="url" class="media-thumb" />
              <div v-else class="media-file-icon">📄</div>
              <span class="media-name">{{ url.split('/').pop() }}</span>
            </div>
          </div>
        </div>
      </el-form>
    </el-card>
  </div>
</template>
