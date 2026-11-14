import { useState } from "react";

const apiUrl = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export default function App() {
  const [message, setMessage] = useState("");

  async function fetchMessage() {
    const response = await fetch(apiUrl);
    const data: { message: string } = await response.json();
    setMessage(data.message);
  }

  return (
    <main className="mx-auto max-w-xl p-8">
      <h1 className="text-2xl font-bold">demo</h1>
      <button className="mt-4 rounded bg-blue-600 px-4 py-2 text-white" onClick={fetchMessage}>
        Fetch message from the API
      </button>
      {message && <p className="mt-4">{message}</p>}
    </main>
  );
}
