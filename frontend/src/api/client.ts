import type { ApiError } from '../types'

class ApiClient {
  private baseUrl = '/api'

  private async request<T>(method: string, path: string, body?: unknown): Promise<T> {
    const url = `${this.baseUrl}${path}`
    const isGet = method === 'GET'

    try {
      const options: RequestInit = {
        method,
        headers: {
          'Content-Type': 'application/json'
        }
      }

      if (body && !isGet) {
        options.body = JSON.stringify(body)
      }

      const response = await fetch(url, options)

      if (response.ok) {
        const data = await response.json() as T
        if (isGet) {
          try {
            localStorage.setItem(`cache:${path}`, JSON.stringify(data))
          } catch {
            // localStorage full, ignore
          }
        }
        return data
      }

      if (response.status >= 400 && response.status < 500) {
        try {
          const err = await response.json() as ApiError
          throw new Error(err.message || `Erro ${response.status}`)
        } catch (e) {
          if (e instanceof Error && e.message !== `Erro ${response.status}`) {
            throw e
          }
          throw new Error(`Erro na requisição: ${response.status}`)
        }
      }

      throw new Error('Erro interno do servidor. Tente novamente.')
    } catch (error) {
      if (error instanceof TypeError) {
        // Network error
        if (isGet) {
          const cached = localStorage.getItem(`cache:${path}`)
          if (cached) {
            return JSON.parse(cached) as T
          }
        }
        throw new Error('Sem conexão com o servidor.')
      }
      throw error
    }
  }

  async get<T>(path: string): Promise<T> {
    return this.request<T>('GET', path)
  }

  async post<T>(path: string, body?: unknown): Promise<T> {
    return this.request<T>('POST', path, body)
  }

  async put<T>(path: string, body?: unknown): Promise<T> {
    return this.request<T>('PUT', path, body)
  }

  async patch<T>(path: string, body?: unknown): Promise<T> {
    return this.request<T>('PATCH', path, body)
  }

  async del(path: string): Promise<void> {
    await this.request<void>('DELETE', path)
  }
}

export const api = new ApiClient()
