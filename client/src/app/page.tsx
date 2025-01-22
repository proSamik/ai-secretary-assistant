'use client';

import {TodoList} from "../components/TodoList";
import React from "react";

export default function Home() {
  return (
    <main className="min-h-screen bg-gray-50 py-8">
      <TodoList />
    </main>
  );
}