'use client';

import { Todo } from '@/lib/api';

interface TodoItemProps {
  todo: Todo;
  onStatusChange: (id: number, status: 'pending' | 'completed') => void;
  onDelete: (id: number) => void;
}

export function TodoItem({ todo, onStatusChange, onDelete }: TodoItemProps) {
  const toggleStatus = () => {
    const newStatus = todo.status === 'pending' ? 'completed' : 'pending';
    onStatusChange(todo.id, newStatus);
  };

  return (
    <div className={`
      bg-white rounded-lg shadow-md p-4 mb-4
      border-l-4 ${todo.status === 'completed' ? 'border-green-500' : 'border-yellow-500'}
      transition-all duration-200 ease-in-out hover:shadow-lg
    `}>
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1 min-w-0">
          <h3 className="text-lg font-semibold text-gray-900 truncate">
            {todo.title}
          </h3>
          <p className="mt-1 text-gray-600 text-sm">
            {todo.description}
          </p>
          <p className="mt-2 text-sm text-gray-500">
            Due: {new Date(todo.due_date).toLocaleDateString()}
          </p>
        </div>
        <div className="flex gap-2 shrink-0">
          <button
            onClick={toggleStatus}
            className={`
              px-3 py-1 rounded-md text-sm font-medium
              transition-colors duration-150 ease-in-out
              ${todo.status === 'completed'
                ? 'bg-yellow-100 text-yellow-700 hover:bg-yellow-200'
                : 'bg-green-100 text-green-700 hover:bg-green-200'
              }
            `}
          >
            {todo.status === 'completed' ? 'Mark Pending' : 'Mark Complete'}
          </button>
          <button
            onClick={() => onDelete(todo.id)}
            className="px-3 py-1 rounded-md text-sm font-medium bg-red-100 text-red-700 hover:bg-red-200 transition-colors duration-150 ease-in-out"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  );
}