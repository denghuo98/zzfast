<!-- 烟花效果：复刻 art-design-pro FireworkSystem，动态创建 canvas 避免 ClientOnly + ref 时序问题 -->
<template>
  <div />
</template>

<script setup lang="ts">
import { useEventListener } from '@vueuse/core'

// ===== 粒子配置 =====
const CONFIG = {
  POOL_SIZE: 600,
  PARTICLES_PER_BURST: 200,
  SIZES: {
    RECTANGLE: { WIDTH: 24, HEIGHT: 12 },
    SQUARE: { SIZE: 12 },
    CIRCLE: { SIZE: 12 },
    TRIANGLE: { SIZE: 10 },
    OVAL: { WIDTH: 24, HEIGHT: 12 },
  },
  ROTATION: { BASE_SPEED: 2, RANDOM_SPEED: 3, DECAY: 0.98 },
  PHYSICS: { GRAVITY: 0.525, VELOCITY_THRESHOLD: 10, OPACITY_DECAY: 0.02 },
  COLORS: [
    'rgba(255,68,68,1)', 'rgba(255,68,68,0.9)', 'rgba(255,68,68,0.8)',
    'rgba(255,116,188,1)', 'rgba(255,116,188,0.9)', 'rgba(255,116,188,0.8)',
    'rgba(68,68,255,0.8)', 'rgba(92,202,56,0.7)', 'rgba(255,68,255,0.8)',
    'rgba(68,255,255,0.7)', 'rgba(255,136,68,0.7)', 'rgba(68,136,255,1)',
    'rgba(250,198,122,0.8)',
  ],
  SHAPES: ['rectangle','rectangle','rectangle','rectangle','rectangle','rectangle','rectangle','circle','triangle','oval'],
} as const

interface Particle {
  x: number; y: number; vx: number; vy: number
  color: string; rotation: number; rotationSpeed: number
  scale: number; shape: string; opacity: number; active: boolean
}

let canvas: HTMLCanvasElement | null = null
let ctx: CanvasRenderingContext2D | null = null
let w = 0; let h = 0
let pool: Particle[] = []
let active: Particle[] = []
let poolIdx = 0
let animId = 0

function createParticle(): Particle {
  return { x:0,y:0,vx:0,vy:0,color:'',rotation:0,rotationSpeed:0,scale:1,shape:'circle',opacity:1,active:false }
}

function getFromPool(): Particle | null {
  for (let i = 0; i < CONFIG.POOL_SIZE; i++) {
    const idx = (poolIdx + i) % CONFIG.POOL_SIZE
    const p = pool[idx]!
    if (!p.active) { poolIdx = (idx + 1) % CONFIG.POOL_SIZE; p.active = true; return p }
  }
  return null
}

function recycle(idx: number) {
  const p = active[idx]!
  p.active = false
  active.splice(idx, 1)
}

function createFirework(cx?: number, cy?: number) {
  const startX = cx ?? Math.random() * w
  const startY = cy ?? h

  for (let i = 0; i < CONFIG.PARTICLES_PER_BURST; i++) {
    const particle = getFromPool()
    if (!particle) continue

    const angle = (Math.PI * i) / (CONFIG.PARTICLES_PER_BURST / 2)
    const speed = (12 + Math.random() * 6) * 1.5
    const spread = Math.random() * Math.PI * 2

    particle.x = startX
    particle.y = startY
    particle.vx = Math.cos(angle) * Math.cos(spread) * speed * (Math.random() * 0.5 + 0.5)
    particle.vy = Math.sin(angle) * speed - 15
    particle.color = CONFIG.COLORS[Math.floor(Math.random() * CONFIG.COLORS.length)]!
    particle.rotation = Math.random() * 360
    particle.rotationSpeed = (Math.random() * CONFIG.ROTATION.RANDOM_SPEED + CONFIG.ROTATION.BASE_SPEED) * (Math.random() > 0.5 ? 1 : -1)
    particle.scale = 0.8 + Math.random() * 0.4
    particle.shape = CONFIG.SHAPES[Math.floor(Math.random() * CONFIG.SHAPES.length)]!
    particle.opacity = 1

    active.push(particle)
  }
}

function renderShape(p: Particle) {
  if (!ctx) return
  const s = CONFIG.SIZES
  ctx.fillStyle = p.color

  switch (p.shape) {
    case 'rectangle':
      ctx.fillRect(-s.RECTANGLE.WIDTH / 2, -s.RECTANGLE.HEIGHT / 2, s.RECTANGLE.WIDTH, s.RECTANGLE.HEIGHT)
      break
    case 'square':
      ctx.fillRect(-s.SQUARE.SIZE / 2, -s.SQUARE.SIZE / 2, s.SQUARE.SIZE, s.SQUARE.SIZE)
      break
    case 'circle':
      ctx.beginPath(); ctx.arc(0, 0, s.CIRCLE.SIZE / 2, 0, Math.PI * 2); ctx.fill(); break
    case 'triangle':
      ctx.beginPath(); ctx.moveTo(0, -s.TRIANGLE.SIZE); ctx.lineTo(s.TRIANGLE.SIZE, s.TRIANGLE.SIZE); ctx.lineTo(-s.TRIANGLE.SIZE, s.TRIANGLE.SIZE); ctx.closePath(); ctx.fill()
      break
    case 'oval':
      ctx.beginPath(); ctx.ellipse(0, 0, s.OVAL.WIDTH / 2, s.OVAL.HEIGHT / 2, 0, 0, Math.PI * 2); ctx.fill(); break
  }
}

function animate() {
  if (!ctx || !canvas) return
  const { GRAVITY, VELOCITY_THRESHOLD, OPACITY_DECAY } = CONFIG.PHYSICS
  const { DECAY } = CONFIG.ROTATION

  ctx.clearRect(0, 0, w, h)
  ctx.globalCompositeOperation = 'lighter'

  for (let i = active.length - 1; i >= 0; i--) {
    const p = active[i]!
    p.x += p.vx; p.y += p.vy; p.vy += GRAVITY
    p.rotation += p.rotationSpeed; p.rotationSpeed *= DECAY

    if (p.vy > VELOCITY_THRESHOLD) {
      p.opacity -= OPACITY_DECAY
      if (p.opacity <= 0) { recycle(i); continue }
    }
    if (p.x < -100 || p.x > w + 100 || p.y < -100 || p.y > h + 100) { recycle(i); continue }

    ctx.save()
    ctx.globalAlpha = p.opacity
    ctx.translate(p.x, p.y)
    ctx.rotate((p.rotation * Math.PI) / 180)
    ctx.scale(p.scale, p.scale)
    renderShape(p)
    ctx.restore()
  }

  animId = requestAnimationFrame(animate)
}

function resize() {
  if (!canvas) return
  canvas.width = w = window.innerWidth
  canvas.height = h = window.innerHeight
}

onMounted(() => {
  if (typeof document === 'undefined') return

  // 初始化对象池
  for (let i = 0; i < CONFIG.POOL_SIZE; i++) pool.push(createParticle())

  // 动态创建 canvas 并插入 body
  canvas = document.createElement('canvas')
  canvas.className = 'fixed top-0 left-0 z-[9999] w-full h-full pointer-events-none'
  document.body.appendChild(canvas)

  ctx = canvas.getContext('2d')
  if (!ctx) return

  resize()

  useEventListener('resize', resize)
  useEventListener('keydown', (e: KeyboardEvent) => {
    if ((e.metaKey || e.ctrlKey) && e.shiftKey && e.key.toLowerCase() === 'p') {
      e.preventDefault()
      createFirework()
    }
  })

  animId = requestAnimationFrame(animate)
})

onBeforeUnmount(() => {
  cancelAnimationFrame(animId)
  canvas?.remove()
  active = []
  pool = []
})

defineExpose({ createFirework })
</script>
