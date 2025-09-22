<script setup lang="ts">
import { RouterView, useRouter } from 'vue-router'
import { watch } from 'vue'

// 获取路由实例
const router = useRouter()

// 监听路由变化时设置过渡名称
let transitionName = 'slide-left'

watch(() => router.currentRoute.value, (to, from) => {
  // 根据路由变化方向设置不同的过渡效果
  if (from && to) {
    // 根据路由路径判断是前进还是后退
    const routesOrder = ['login', 'register', 'roleSelect', 'chat']
    const fromIndex = routesOrder.indexOf(from.name as string)
    const toIndex = routesOrder.indexOf(to.name as string)

    // 如果找到路由索引，根据索引变化判断方向
    if (fromIndex !== -1 && toIndex !== -1) {
      transitionName = toIndex > fromIndex ? 'slide-left' : 'slide-right'
    } else {
      // 默认使用淡入淡出效果
      transitionName = 'fade'
    }
  }
})
</script>

<template>
  <Transition :name="transitionName" mode="out-in">
    <RouterView />
  </Transition>
</template>

<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background-color: #f5f5f5;
}

/* 淡入淡出过渡效果 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 左滑过渡效果 - 用于前进导航 */
.slide-left-enter-active,
.slide-left-leave-active {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.slide-left-enter-from {
  transform: translateX(100%);
  opacity: 0;
}

.slide-left-leave-to {
  transform: translateX(-100%);
  opacity: 0;
}

/* 右滑过渡效果 - 用于后退导航 */
.slide-right-enter-active,
.slide-right-leave-active {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.slide-right-enter-from {
  transform: translateX(-100%);
  opacity: 0;
}

.slide-right-leave-to {
  transform: translateX(100%);
  opacity: 0;
}

/* 确保路由视图充满整个屏幕 */
.router-view {
  min-height: 100vh;
  width: 100%;
}
</style>
