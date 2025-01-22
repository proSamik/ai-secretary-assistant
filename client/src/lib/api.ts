const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

export interface Todo {
  id: number;
  title: string;
  description: string;
  status: 'pending' | 'completed';
  due_date: string;
  created_at: string;
  updated_at: string;
}

export interface TodoInput {
  title: string;
  description: string;
  due_date: string;
}

export const api = {
  async getTodos() {
    const response = await fetch(`${API_BASE}/todos`);
    if (!response.ok) throw new Error('Failed to fetch todos');
    return response.json() as Promise<Todo[]>;
  },

  async createTodo(todo: TodoInput) {
    const response = await fetch(`${API_BASE}/todos`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(todo),
    });
    if (!response.ok) throw new Error('Failed to create todo');
    return response.json() as Promise<Todo>;
  },

  async updateTodoStatus(id: number, status: 'pending' | 'completed') {
    const response = await fetch(`${API_BASE}/todos/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ status }),
    });
    if (!response.ok) throw new Error('Failed to update todo');
    return response.json() as Promise<Todo>;
  },

  async deleteTodo(id: number) {
    const response = await fetch(`${API_BASE}/todos/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) throw new Error('Failed to delete todo');
    return true;
  },

  async getTodo(id: number) {
    const response = await fetch(`${API_BASE}/todos/${id}`);
    if (!response.ok) throw new Error('Failed to fetch todo');
    return response.json() as Promise<Todo>;
  },
};