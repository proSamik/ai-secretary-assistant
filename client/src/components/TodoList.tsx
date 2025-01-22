'use client';

import React, { useState, useEffect } from 'react';
import { TodoItem } from './TodoItem';
import { TodoForm } from './TodoForm';
import {api, Todo} from "../lib/api";

export function TodoList() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchTodos = async () => {
    try {
      const data = await api.getTodos();
      setTodos(data);
      setError(null);
    } catch (err) {
      setError('Failed to load todos');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTodos();
  }, []);

  const handleAddTodo = async (todo: Todo) => {
    setTodos(prev => [todo, ...prev]);
  };

  const handleUpdateTodo = async (id: number, status: 'pending' | 'completed') => {
    try {
      await api.updateTodoStatus(id, status);
      fetchTodos();
    } catch (err) {
      setError('Failed to update todo');
    }
  };

  const handleDeleteTodo = async (id: number) => {
    try {
      await api.deleteTodo(id);
      setTodos(prev => prev.filter(todo => todo.id !== id));
    } catch (err) {
      setError('Failed to delete todo');
    }
  };

  if (loading) {
    return <div className="flex justify-center p-4">Loading...</div>;
  }

  return (
    <div className="container mx-auto p-4 max-w-4xl">
      <h1 className="text-3xl font-bold mb-8">My Todos</h1>
      
      <TodoForm onAddTodo={handleAddTodo} />

      {error && (
        <div className="bg-red-100 text-red-700 p-4 rounded mb-4">
          {error}
        </div>
      )}

      <div className="space-y-4">
        {todos.map((todo) => (
          <TodoItem
            key={todo.id}
            todo={todo}
            onStatusChange={handleUpdateTodo}
            onDelete={handleDeleteTodo}
          />
        ))}
        {todos.length === 0 && (
          <div className="text-center text-gray-500 p-8">
            No todos yet. Add one above!
          </div>
        )}
      </div>
    </div>
  );
}