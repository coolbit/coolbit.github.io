<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { ElNotification } from 'element-plus'
import { Plus, Notebook, Upload } from '@element-plus/icons-vue'
import { publish } from '../api/publish'

const route = useRoute()
const postsActive = computed(() => route.path === '/' || route.path.startsWith('/posts'))

const publishing = ref(false)

async function onPublish() {
  publishing.value = true
  try {
    const { data } = await publish()
    ElNotification({
      title: `${data.count} posts exported`,
      message: `<code style="font-size:12px">${data.dir}</code><br>Run <code>git push</code> to deploy.`,
      type: 'success',
      dangerouslyUseHTMLString: true,
      duration: 0,
    })
  } catch (e: any) {
    const msg = e?.response?.data?.error || e?.message || 'Unknown error'
    ElNotification({ title: 'Publish failed', message: msg, type: 'error', duration: 0 })
  } finally {
    publishing.value = false
  }
}
</script>

<template>
  <div>
    <header class="site-header">
      <div class="header-inner">
        <RouterLink to="/" class="site-brand">
          <el-icon :size="20"><Notebook /></el-icon>
          Blog Manager
        </RouterLink>
        <div class="header-right">
          <RouterLink to="/" :class="['nav-link', { active: postsActive }]">Posts</RouterLink>
          <el-button
            :icon="Upload"
            :loading="publishing"
            size="small"
            @click="onPublish"
          >
            {{ publishing ? 'Publishing…' : 'Publish' }}
          </el-button>
          <RouterLink to="/posts/new">
            <el-button type="primary" :icon="Plus">New Post</el-button>
          </RouterLink>
        </div>
      </div>
    </header>
    <main class="main-content">
      <RouterView />
    </main>
  </div>
</template>
