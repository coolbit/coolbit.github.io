<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { EditPen, Delete, Search } from '@element-plus/icons-vue'
import { listPosts, deletePost } from '../api/posts'
import type { Post } from '../types/post'

const router = useRouter()
const posts = ref<Post[]>([])
const loading = ref(true)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')

async function load() {
  loading.value = true
  try {
    const { data } = await listPosts({
      page: page.value,
      page_size: pageSize.value,
      q: searchQuery.value.trim() || undefined,
    })
    posts.value = data.data
    total.value = data.total
  } finally {
    loading.value = false
  }
}

// Debounce search — reset to page 1 on new query
let debounceTimer: ReturnType<typeof setTimeout>
watch(searchQuery, () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    page.value = 1
    load()
  }, 350)
})

async function remove(post: Post) {
  await deletePost(post.id)
  ElMessage.success('Post deleted')
  if (posts.value.length === 1 && page.value > 1) page.value--
  load()
}

function onPaginationChange() {
  load()
}

function fmt(s: string) {
  return new Date(s).toLocaleDateString('zh-CN', { year: 'numeric', month: 'short', day: 'numeric' })
}

onMounted(load)
</script>

<template>
  <div>
    <div class="page-header">
      <h2>
        Posts
        <el-tag type="info" size="small">{{ total }}</el-tag>
      </h2>
      <el-input
        v-model="searchQuery"
        placeholder="Search title, summary, content…"
        :prefix-icon="Search"
        clearable
        style="width: 300px"
      />
    </div>

    <el-table :data="posts" v-loading="loading" stripe style="width: 100%; border-radius: 8px" table-layout="auto">
      <template #empty>
        <el-empty
          :description="searchQuery ? `No results for &quot;${searchQuery}&quot;` : 'No posts yet'"
          :image-size="120"
        >
          <el-button v-if="!searchQuery" type="primary" @click="router.push('/posts/new')">
            Write your first post
          </el-button>
        </el-empty>
      </template>

      <el-table-column label="Title" min-width="220">
        <template #default="{ row }">
          <RouterLink :to="`/posts/${row.id}`" class="post-link">{{ row.title }}</RouterLink>
        </template>
      </el-table-column>

      <el-table-column label="Category" width="150">
        <template #default="{ row }">
          <el-tag v-if="row.category" size="small" effect="light" round>{{ row.category }}</el-tag>
          <span v-else class="text-muted">—</span>
        </template>
      </el-table-column>

      <el-table-column label="Summary" min-width="200" show-overflow-tooltip>
        <template #default="{ row }">
          <span class="text-muted">{{ row.summary || '—' }}</span>
        </template>
      </el-table-column>

      <el-table-column label="Date" width="130">
        <template #default="{ row }">
          <span class="text-muted">{{ fmt(row.created_at) }}</span>
        </template>
      </el-table-column>

      <el-table-column label="Actions" width="180" align="right">
        <template #default="{ row }">
          <el-button :icon="EditPen" size="small" @click="router.push(`/posts/${row.id}/edit`)">
            Edit
          </el-button>
          <el-popconfirm
            :title="`Delete &quot;${row.title}&quot;?`"
            confirm-button-text="Delete"
            confirm-button-type="danger"
            cancel-button-text="Cancel"
            width="220"
            @confirm="remove(row)"
          >
            <template #reference>
              <el-button :icon="Delete" size="small" type="danger" plain>Delete</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <div style="display: flex; justify-content: flex-end; margin-top: 1.25rem">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        background
        @change="onPaginationChange"
      />
    </div>
  </div>
</template>
