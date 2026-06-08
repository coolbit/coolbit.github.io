<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Check, ArrowLeft } from '@element-plus/icons-vue'
import { MdEditor } from 'md-editor-v3'
import { getPost, createPost, updatePost } from '../api/posts'
import { listCategories } from '../api/categories'

const route = useRoute()
const router = useRouter()

const postId = computed(() => route.params.id ? Number(route.params.id) : null)
const isEdit = computed(() => postId.value !== null)

const title = ref('')
const summary = ref('')
const content = ref('')
const category = ref('')
const categories = ref<string[]>([])
const saving = ref(false)

onMounted(async () => {
  const [catsRes] = await Promise.all([
    listCategories(),
    isEdit.value
      ? getPost(postId.value!).then(({ data }) => {
          title.value = data.title
          summary.value = data.summary
          content.value = data.content
          category.value = data.category
        })
      : Promise.resolve(),
  ])
  categories.value = catsRes.data
})

async function save() {
  if (!title.value.trim()) {
    ElMessage.warning('Title is required')
    return
  }
  saving.value = true
  try {
    const payload = {
      title: title.value.trim(),
      summary: summary.value.trim(),
      content: content.value,
      category: category.value,
    }
    if (isEdit.value) {
      await updatePost(postId.value!, payload)
      ElMessage.success('Post saved')
      router.push(`/posts/${postId.value}`)
    } else {
      const { data } = await createPost(payload)
      ElMessage.success('Post created')
      router.push(`/posts/${data.id}`)
    }
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div>
    <div class="page-header">
      <h2>{{ isEdit ? 'Edit Post' : 'New Post' }}</h2>
      <div style="display: flex; gap: 8px">
        <el-button :icon="ArrowLeft" @click="router.back()">Cancel</el-button>
        <el-button type="primary" :icon="Check" :loading="saving" @click="save">Save</el-button>
      </div>
    </div>

    <el-card shadow="never" style="border-radius: 10px">
      <el-form label-position="top" class="editor-form">
        <el-row :gutter="16">
          <el-col :xs="24" :md="16">
            <el-form-item label="Title" required>
              <el-input v-model="title" placeholder="Post title" size="large" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="8">
            <el-form-item label="Category">
              <el-select
                v-model="category"
                placeholder="No category"
                clearable
                filterable
                allow-create
                style="width: 100%"
                size="large"
              >
                <el-option v-for="cat in categories" :key="cat" :label="cat" :value="cat" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="Summary">
          <el-input v-model="summary" placeholder="Brief summary (optional)" />
        </el-form-item>

        <el-form-item label="Content">
          <MdEditor v-model="content" style="width: 100%; height: 520px" />
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>
