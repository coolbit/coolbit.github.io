<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MdPreview } from 'md-editor-v3'
import { ArrowLeft, EditPen } from '@element-plus/icons-vue'
import { getPost } from '../api/posts'
import type { Post } from '../types/post'

const route = useRoute()
const router = useRouter()
const post = ref<Post | null>(null)

onMounted(async () => {
  const { data } = await getPost(Number(route.params.id))
  post.value = data
})

function fmt(s: string) {
  return new Date(s).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })
}
</script>

<template>
  <div>
    <div v-if="!post" v-loading="true" style="min-height: 200px" element-loading-text="Loading…" />

    <template v-else>
      <div class="page-header">
        <el-button :icon="ArrowLeft" @click="router.back()">Back</el-button>
        <el-button :icon="EditPen" @click="router.push(`/posts/${post.id}/edit`)">Edit</el-button>
      </div>

      <el-card shadow="never" style="border-radius: 10px">
        <article class="post-detail">
          <div class="post-meta-row">
            <span class="post-date">{{ fmt(post.created_at) }}</span>
          </div>

          <h1 class="post-title">{{ post.title }}</h1>

          <p v-if="post.summary" class="post-summary">{{ post.summary }}</p>

          <el-divider />

          <div class="post-content">
            <MdPreview :modelValue="post.content" />
          </div>
        </article>
      </el-card>
    </template>
  </div>
</template>
