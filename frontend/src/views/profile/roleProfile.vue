<template>
  <div class="role-profile-container" :style="backgroundStyle">
    <!-- 头部导航 -->
    <div class="top-navigation">
      <div class="logo-container">
        <button class="back-button" @click="goBack">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M19 12H5M5 12L12 5M5 12L12 19" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
        <div class="logo">
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" fill="rgba(255,255,255,0.9)"/>
            <path d="M2 17L12 22L22 17" stroke="rgba(255,255,255,0.9)" stroke-width="2" stroke-linecap="round"/>
            <path d="M2 12L12 17L22 12" stroke="rgba(255,255,255,0.9)" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <span class="logo-text">角色详情</span>
        </div>
      </div>
    </div>

    <!-- 主内容区域 -->
    <div class="role-profile-content">
      <!-- 角色基本信息 -->
      <div class="role-basic-info">
        <div class="avatar-container">
          <div
            class="role-avatar"
            :style="{ backgroundImage: `url(${roleInfo.avatarUrl})` }"
          ></div>
        </div>
        <div class="role-details">
          <h1 class="role-name">{{ roleInfo.name }}</h1>
          <div class="role-category">{{ roleInfo.category }}</div>
          <div class="role-stats">
            <div class="stat-item">
              <span class="stat-label">对话次数</span>
              <span class="stat-value">{{ roleStats.conversationCount }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">收藏人数</span>
              <span class="stat-value">{{ roleStats.favoriteCount }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">评分</span>
              <span class="stat-value">{{ roleStats.rating }}</span>
            </div>
          </div>
        </div>

        <!-- 右侧操作按钮区域 -->
        <div class="action-buttons">
          <button class="start-chat-btn" @click="startChat">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            开始对话
          </button>
          <button
            class="favorite-btn"
            :class="{ favorited: isFavorite }"
            @click="toggleFavorite"
          >
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" :fill="isFavorite ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            {{ isFavorite ? '已收藏' : '收藏' }}
          </button>
        </div>
      </div>

      <!-- 标签导航 -->
      <div class="tab-navigation">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          class="tab-button"
          :class="{ active: activeTab === tab.id }"
          @click="activeTab = tab.id"
        >
          {{ tab.name }}
        </button>
      </div>

      <!-- 内容区域 -->
      <div class="content-area">
        <!-- 角色介绍 -->
        <div v-if="activeTab === 'introduction'" class="tab-content">
          <div class="content-card">
            <h2 class="section-title">角色介绍</h2>
            <div class="content-text">{{ roleInfo.fullDescription }}</div>
          </div>

          <div class="content-card">
            <h2 class="section-title">基本信息</h2>
            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">所属时代</span>
                <span class="info-value">{{ roleInfo.era }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">出生地</span>
                <span class="info-value">{{ roleInfo.birthplace }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">代表作品</span>
                <span class="info-value">{{ roleInfo.representativeWorks }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">特点</span>
                <span class="info-value">{{ roleInfo.characteristics.join('、') }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 历史背景 -->
        <div v-if="activeTab === 'background'" class="tab-content">
          <div class="content-card">
            <h2 class="section-title">历史背景</h2>
            <div class="content-text">{{ roleInfo.historicalBackground }}</div>
          </div>

          <div class="content-card timeline">
            <h2 class="section-title">重要时间线</h2>
            <div v-for="event in roleInfo.timelineEvents" :key="event.year" class="timeline-item">
              <div class="timeline-year">{{ event.year }}</div>
              <div class="timeline-content">
                <h3 class="timeline-title">{{ event.title }}</h3>
                <p class="timeline-description">{{ event.description }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 主要事迹 -->
        <div v-if="activeTab === 'achievements'" class="tab-content">
          <div v-for="achievement in roleInfo.achievements" :key="achievement.id" class="content-card">
            <h2 class="section-title">{{ achievement.title }}</h2>
            <div class="achievement-meta">
              <span class="achievement-time">{{ achievement.time }}</span>
            </div>
            <div class="content-text">{{ achievement.description }}</div>
          </div>
        </div>

        <!-- 行事风格 -->
        <div v-if="activeTab === 'style'" class="tab-content">
          <div class="content-card">
            <h2 class="section-title">行事风格</h2>
            <div class="content-text">{{ roleInfo.behaviorStyle }}</div>
          </div>

          <div class="content-card">
            <h2 class="section-title">经典语句</h2>
            <div class="quotes-container">
              <div v-for="quote in roleInfo.famousQuotes" :key="quote.id" class="quote-card">
                <svg class="quote-icon" width="32" height="32" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M10 11H7C5.89543 11 5 11.8954 5 13V17C5 18.1046 5.89543 19 7 19H10V11Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  <path d="M17 11H14C12.8954 11 12 11.8954 12 13V17C12 18.1046 12.8954 19 14 19H17V11Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                <p class="quote-text">{{ quote.text }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import type { Role } from '@/types/api';

const router = useRouter();
const route = useRoute();
const roleId = route.params.roleId as string;

// 标签导航
const tabs = [
  { id: 'introduction', name: '角色介绍' },
  { id: 'background', name: '历史背景' },
  { id: 'achievements', name: '主要事迹' },
  { id: 'style', name: '行事风格' }
];
const activeTab = ref('introduction');

// 角色信息（扩展了基本Role类型，添加了更多详情字段）
interface ExtendedRole extends Role {
  fullDescription: string;
  era: string;
  birthplace: string;
  representativeWorks: string;
  characteristics: string[];
  historicalBackground: string;
  timelineEvents: Array<{ year: string; title: string; description: string }>;
  achievements: Array<{ id: string; title: string; time: string; description: string }>;
  behaviorStyle: string;
  famousQuotes: Array<{ id: string; text: string }>;
  colorScheme: { primary: string; secondary: string; accent: string };
}

// 角色统计信息
interface RoleStats {
  conversationCount: number;
  favoriteCount: number;
  rating: number;
}

// 角色信息（模拟数据，实际项目中应从API获取）
const roleInfo = ref<ExtendedRole>({
  id: roleId || '1',
  name: '夏洛克·福尔摩斯',
  category: '文学角色',
  avatarUrl: 'https://picsum.photos/id/1/400/400',
  description: '英国著名侦探，善于观察和推理',
  fullDescription: '夏洛克·福尔摩斯（Sherlock Holmes）是英国作家阿瑟·柯南·道尔创作的虚构侦探，以其卓越的观察力、推理能力和逻辑思维而闻名。他善于通过细微的线索解决复杂的案件，常使用演绎推理法来破解谜团。福尔摩斯居住在伦敦贝克街221B号，与他的好友兼助手约翰·华生医生一起经历了许多惊险刺激的冒险。他冷静、机智，对细节有着惊人的洞察力，同时也有些特立独行，喜欢拉小提琴、做化学实验，并且对某些领域有着深入的研究。',
  era: '19世纪末至20世纪初',
  birthplace: '英国',
  representativeWorks: '《血字的研究》、《四签名》、《巴斯克维尔的猎犬》',
  characteristics: ['观察力敏锐', '逻辑推理能力强', '冷静理性', '特立独行', '精通化学', '擅长小提琴'],
  historicalBackground: '夏洛克·福尔摩斯的故事背景设定在19世纪末的维多利亚时代伦敦。这一时期是英国工业革命的高峰期，伦敦作为大英帝国的首都，既呈现出繁荣的一面，也存在着严重的社会问题。贫富差距悬殊，犯罪率居高不下，这为侦探故事提供了丰富的素材。同时，这一时期科学技术和法医学开始发展，为福尔摩斯使用科学方法破案提供了背景支持。',
  timelineEvents: [
    { year: '1887年', title: '首次登场', description: '在阿瑟·柯南·道尔的小说《血字的研究》中首次出现' },
    { year: '1891年', title: '与莫里亚蒂教授对决', description: '在《最后一案》中，福尔摩斯与宿敌莫里亚蒂教授在莱辛巴赫瀑布决斗' },
    { year: '1894年', title: '归来记', description: '在读者的强烈要求下，道尔在《空屋》中让福尔摩斯重返舞台' },
    { year: '1902年', title: '巴斯克维尔的猎犬', description: '发表《巴斯克维尔的猎犬》，这是福尔摩斯最著名的案件之一' }
  ],
  achievements: [
    { id: '1', title: '破解血字的研究', time: '1881年', description: '福尔摩斯通过分析犯罪现场的血字和其他线索，成功破获了一起复仇杀人案，展现了他卓越的推理能力。' },
    { id: '2', title: '四签名奇案', time: '1888年', description: '福尔摩斯破解了涉及印度宝藏和四个签名的复杂案件，揭示了隐藏多年的秘密。' },
    { id: '3', title: '巴斯克维尔的猎犬', time: '1890年', description: '福尔摩斯调查了巴斯克维尔家族的诅咒传说，揭露了利用传说进行谋杀的阴谋。' },
    { id: '4', title: '艾琳·艾德勒事件', time: '1888年', description: '在《波西米亚丑闻》中，福尔摩斯遇到了他一生中最尊敬的对手艾琳·艾德勒，虽然未能完全成功完成任务，但对他评价极高。' }
  ],
  behaviorStyle: '福尔摩斯的行事风格冷静、理性，注重逻辑和证据。他总是以客观的态度分析案件，不受情感因素干扰。他善于从看似无关的细节中发现重要线索，常说"当你排除一切不可能的情况，剩下的，不管多难以置信，那都是事实"。福尔摩斯喜欢独来独往，但与华生医生保持着深深的友谊。他有时显得有些冷漠和不近人情，但内心深处有着强烈的正义感，总是站在正义的一方。在调查案件时，他经常乔装改扮，深入虎穴，展现出非凡的勇气和智慧。',
  famousQuotes: [
    { id: '1', text: '当你排除一切不可能的情况，剩下的，不管多难以置信，那都是事实。' },
    { id: '2', text: '我从不猜测，猜测是一个很坏的习惯，会影响正常的逻辑推理能力。' },
    { id: '3', text: '对于一个真正的推理家而言，如果有人指给他一个事实的其中一个方面，他不仅能推断出这个事实的各个方面，而且能够推断出由此将会产生的一切后果。' },
    { id: '4', text: '在没有得到任何证据的情况下是不能进行推理的，那样的话，只能是误入歧途。' }
  ],
  colorScheme: { primary: '#1a237e', secondary: '#5c6bc0', accent: '#8e24aa' }
});

// 角色统计信息
const roleStats = ref<RoleStats>({ conversationCount: 1254, favoriteCount: 892, rating: 4.8 });

// 是否收藏
const isFavorite = ref(false);

// 计算背景样式
const backgroundStyle = computed(() => {
  const colors = roleInfo.value.colorScheme;
  return { background: `linear-gradient(135deg, ${colors.primary} 0%, ${colors.secondary} 100%)`, backgroundAttachment: 'fixed' };
});

// 返回上一页
const goBack = () => router.back();

// 开始对话
const startChat = () => router.push({ name: 'chat', params: { roleId: roleInfo.value.id } });

// 切换收藏状态
const toggleFavorite = () => {
  isFavorite.value = !isFavorite.value;
  console.log(`${roleInfo.value.name} 收藏状态: ${isFavorite.value ? '已收藏' : '未收藏'}`);
};

// 组件挂载时的初始化
onMounted(() => {
  console.log(`加载角色详情: ${roleId}`);
  isFavorite.value = Math.random() > 0.5; // 模拟随机收藏状态
});
</script>

<style scoped>
/* 角色详情页面容器 */
.role-profile-container {
  min-height: 100vh;
  color: white;
  padding: 1rem;
  background-size: cover;
  background-position: center;
  position: relative;
  overflow: hidden;
}

/* 背景装饰 */
.role-profile-container::before {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(circle at 50% 50%, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0) 70%);
  pointer-events: none;
}

/* 头部导航 */
.top-navigation {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  position: relative;
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.back-button {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.back-button:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateX(-5px);
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.logo-text {
  color: white;
  font-size: 1.5rem;
  font-weight: 700;
  letter-spacing: 0.5px;
}

/* 主内容区域 */
.role-profile-content {
  max-width: 1200px;
  margin: 0 auto;
  position: relative;
}

/* 角色基本信息 */
.role-basic-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 2rem;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 2rem;
  margin-bottom: 2rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.avatar-container {
  position: relative;
}

.role-avatar {
  width: 200px;
  height: 200px;
  border-radius: 50%;
  background-size: cover;
  background-position: center;
  border: 5px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
}

/* 右侧操作按钮区域 */
.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  min-width: 180px;
}

.start-chat-btn, .favorite-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 25px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.start-chat-btn:hover {
  background: white;
  color: #667eea;
  transform: translateY(-3px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
}

.favorite-btn:hover {
  background: rgba(255, 193, 7, 0.2);
  border-color: rgba(255, 193, 7, 0.5);
  transform: translateY(-3px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
}

.favorite-btn.favorited {
  color: #ffc107;
  border-color: #ffc107;
  background: rgba(255, 193, 7, 0.1);
}

.role-details {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 1rem;
  min-width: 0;
}

.role-name {
  font-size: 2.5rem;
  font-weight: 800;
  margin: 0;
  background: linear-gradient(90deg, white, rgba(255, 255, 255, 0.7));
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.role-category {
  font-size: 1.1rem;
  color: rgba(255, 255, 255, 0.8);
  font-style: italic;
  opacity: 0.9;
  margin-bottom: 0.5rem;
}

.role-stats {
  display: flex;
  gap: 2rem;
  margin-top: 0.5rem;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25rem;
}

.stat-label {
  font-size: 0.9rem;
  color: rgba(255, 255, 255, 0.7);
  opacity: 0.8;
  margin-bottom: 0.25rem;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: white;
}

/* 标签导航 */
.tab-navigation {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 2rem;
  overflow-x: auto;
  padding-bottom: 0.5rem;
}

.tab-button {
  padding: 0.75rem 1.5rem;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 25px;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 1rem;
  font-weight: 500;
  white-space: nowrap;
}

.tab-button:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.tab-button.active {
  background: white;
  color: #667eea;
  border-color: white;
}

/* 内容区域 */
.content-area, .tab-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.content-card {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
}

.content-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.section-title {
  font-size: 1.75rem;
  font-weight: 700;
  margin: 0 0 1.5rem 0;
  padding-bottom: 0.75rem;
  border-bottom: 2px solid rgba(255, 255, 255, 0.2);
  color: white;
}

.content-text {
  font-size: 1.1rem;
  line-height: 1.8;
  color: rgba(255, 255, 255, 0.9);
}

/* 基本信息网格 */
.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1.5rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-label {
  font-size: 0.9rem;
  color: rgba(255, 255, 255, 0.7);
  font-weight: 500;
}

.info-value {
  font-size: 1rem;
  color: white;
  background: rgba(255, 255, 255, 0.05);
  padding: 0.75rem;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

/* 时间线 */
.timeline {
  position: relative;
  padding-left: 2rem;
}

.timeline::before {
  content: '';
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  width: 3px;
  background: linear-gradient(to bottom, #667eea, #764ba2);
  border-radius: 3px;
}

.timeline-item {
  display: flex;
  gap: 2rem;
  margin-bottom: 3rem;
  position: relative;
  padding-left: 2rem;
  transition: transform 0.3s ease;
}

.timeline-item:hover {
  transform: translateX(5px);
}

.timeline-item:last-child {
  margin-bottom: 0;
}

.timeline-year {
  width: 80px;
  flex-shrink: 0;
  font-size: 1.3rem;
  font-weight: 700;
  color: white;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  padding: 0.75rem 1rem;
  text-align: center;
  position: relative;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.timeline-year::after {
  content: '';
  position: absolute;
  top: 50%;
  right: -2rem;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: white;
  border: 3px solid #667eea;
  transform: translateY(-50%);
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.2);
}

.timeline-content {
  flex: 1;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 1.5rem;
  border: 1px solid rgba(255, 255, 255, 0.1);
  position: relative;
  overflow: hidden;
}

.timeline-content::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(to right, #667eea, #764ba2);
}

.timeline-title {
  font-size: 1.4rem;
  font-weight: 600;
  margin: 0 0 1rem 0;
  color: white;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.timeline-description {
  font-size: 1.05rem;
  line-height: 1.7;
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .timeline {
    padding-left: 1.5rem;
  }

  .timeline-item {
    flex-direction: column;
    gap: 1rem;
    padding-left: 0;
  }

  .timeline-year {
    width: auto;
    align-self: flex-start;
  }

  .timeline-year::after {
    right: 50%;
    bottom: -2rem;
    top: auto;
    transform: translateX(50%);
  }

  .timeline::before {
    height: calc(100% + 2rem);
  }
}

@media (max-width: 1024px) {
  .role-basic-info {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .role-stats {
    justify-content: center;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .role-profile-container {
    padding: 1rem;
  }

  .role-avatar {
    width: 150px;
    height: 150px;
  }

  .role-name {
    font-size: 2rem;
  }

  .action-buttons {
    width: 100%;
    flex-direction: row;
    justify-content: center;
    flex-wrap: wrap;
  }

  .start-chat-btn, .favorite-btn {
    flex: 1;
    min-width: 150px;
  }

  .timeline::before {
    left: 30px;
  }

  .timeline-year {
    width: 60px;
    font-size: 1rem;
  }

  .timeline-item {
    gap: 1rem;
  }

  .quotes-container {
    grid-template-columns: 1fr;
  }

  .content-card {
    padding: 1.5rem;
  }
}

@media (max-width: 480px) {
  .role-name {
    font-size: 1.75rem;
  }

  .role-stats {
    flex-direction: column;
    gap: 1rem;
  }

  .tab-button {
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
  }
}
</style>
