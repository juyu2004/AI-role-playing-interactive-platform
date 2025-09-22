import { useEffect, useMemo, useRef, useState } from 'react'
import './App.css'

type Role = {
  id: string
  name: string
  category: string
  avatarUrl: string
  description: string
}

export default function App() {
  const [roles, setRoles] = useState<Role[]>([])
  const [selectedRole, setSelectedRole] = useState<string>('')
  const [messages, setMessages] = useState<{ sender: 'user' | 'role'; text: string }[]>([])
  const [input, setInput] = useState('')
  const mediaRecorderRef = useRef<MediaRecorder | null>(null)
  const [recording, setRecording] = useState(false)
  const chunksRef = useRef<Blob[]>([])

  useEffect(() => {
    fetch('/api/roles')
      .then((r) => r.json())
      .then((data) => setRoles(data))
      .catch(() => { })
  }, [])

  const canSend = useMemo(() => selectedRole && input.trim().length > 0, [selectedRole, input])

  async function sendText() {
    if (!canSend) return
    const text = input
    setMessages((m) => [...m, { sender: 'user', text }])
    setInput('')
    const resp = await fetch('/api/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ roleId: selectedRole, text }),
    })
    const data = await resp.json()
    setMessages((m) => [...m, { sender: 'role', text: data.text }])
  }

  async function startRecording() {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    const mediaRecorder = new MediaRecorder(stream)
    mediaRecorderRef.current = mediaRecorder
    chunksRef.current = []
    mediaRecorder.ondataavailable = (e) => {
      if (e.data.size > 0) chunksRef.current.push(e.data)
    }
    mediaRecorder.onstop = async () => {
      const blob = new Blob(chunksRef.current, { type: 'audio/webm' })
      console.log('Recorded audio blob', blob)
    }
    mediaRecorder.start()
    setRecording(true)
  }

  function stopRecording() {
    mediaRecorderRef.current?.stop()
    setRecording(false)
  }

  return (
    <div className="container">
      <h1>AI 角色扮演互动平台</h1>

      <section>
        <h2>选择角色</h2>
        <select value={selectedRole} onChange={(e) => setSelectedRole(e.target.value)}>
          <option value="">请选择角色</option>
          {roles.map((r) => (
            <option key={r.id} value={r.id}>
              {r.name} ({r.category})
            </option>
          ))}
        </select>
      </section>

      <section>
        <h2>聊天</h2>
        <div className="chat-box">
          {messages.map((m, i) => (
            <div key={i} className={m.sender === 'user' ? 'msg user' : 'msg role'}>
              <b>{m.sender === 'user' ? '我' : '角色'}:</b> {m.text}
            </div>
          ))}
        </div>
        <div className="controls">
          <input value={input} onChange={(e) => setInput(e.target.value)} placeholder="输入消息" />
          <button disabled={!canSend} onClick={sendText}>发送</button>
          {!recording ? (
            <button onClick={startRecording}>开始录音</button>
          ) : (
            <button onClick={stopRecording}>停止录音</button>
          )}
        </div>
      </section>
    </div>
  )
}

