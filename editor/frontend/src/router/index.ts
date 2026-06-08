import { createRouter, createWebHistory } from 'vue-router'

export default createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('../layouts/CmsLayout.vue'),
      children: [
        { path: '', component: () => import('../views/BlogList.vue') },
        { path: 'posts/new', component: () => import('../views/BlogEditor.vue') },
        { path: 'posts/:id/edit', component: () => import('../views/BlogEditor.vue') },
        { path: 'posts/:id', component: () => import('../views/BlogDetail.vue') },
      ],
    },
  ],
})
