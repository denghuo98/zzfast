/**
 * 前端密码加密工具 —— 与 zzfast 后端 AES-GCM 加解密对齐。
 *
 * 后端约定：明文密码经 AES-GCM 加密后，输出为
 *   Base64(nonce(12 字节) + ciphertext + tag)
 * 后端 verifyPassword 会先 Base64 解码，再 AesGCMDecrypt 还原明文。
 */

/** 将 Uint8Array 编码为 Base64（兼容 Node SSR 与浏览器） */
function toBase64(bytes: Uint8Array): string {
  if (typeof Buffer !== 'undefined') {
    return Buffer.from(bytes).toString('base64')
  }
  let binary = ''
  bytes.forEach(b => (binary += String.fromCharCode(b)))
  return btoa(binary)
}

/**
 * AES-GCM 加密，返回 Base64(nonce + ciphertext + tag)。
 * key 长度必须为 16/24/32 字节（对应 AES-128/192/256-GCM）。
 */
export async function aesGcmEncrypt(plaintext: string, key: string): Promise<string> {
  const keyBytes = new TextEncoder().encode(key)
  const cryptoKey = await crypto.subtle.importKey(
    'raw',
    keyBytes,
    { name: 'AES-GCM' },
    false,
    ['encrypt'],
  )

  const nonce = crypto.getRandomValues(new Uint8Array(12))
  const data = new TextEncoder().encode(plaintext)
  const encrypted = await crypto.subtle.encrypt({ name: 'AES-GCM', iv: nonce }, cryptoKey, data)

  const ciphertext = new Uint8Array(encrypted)
  const result = new Uint8Array(nonce.length + ciphertext.length)
  result.set(nonce, 0)
  result.set(ciphertext, nonce.length)

  return toBase64(result)
}
