import React, { useState } from "react";

type Props = {
  onConnected: (info: { user: string; pass: string; host: string }, dbs: string[]) => void;
};

const ConnectPage: React.FC<Props> = ({ onConnected }) => {
  const [user, setUser] = useState("");
  const [pass, setPass] = useState("");
  const [host, setHost] = useState("");
  const [error, setError] = useState("");

  const handleConnect = async () => {
    setError("");

    const res = await fetch("http://localhost:8080/connect", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ user, pass, host }),
    });

    if (!res.ok) {
      const t = await res.text();
      setError(t);
      return;
    }

    const data = await res.json();

    onConnected({ user, pass, host }, data.databases);
  };

  return (
    <div style={{ maxWidth: 500, margin: "40px auto" }}>
      <h2>MySQL 接続</h2>

      <input
        placeholder="ユーザー名"
        value={user}
        onChange={(e) => setUser(e.target.value)}
        style={{ width: "100%", marginBottom: 10 }}
      />

      <input
        placeholder="パスワード"
        type="password"
        value={pass}
        onChange={(e) => setPass(e.target.value)}
        style={{ width: "100%", marginBottom: 10 }}
      />

      <input
        placeholder="ホスト (例: localhost:3306)"
        value={host}
        onChange={(e) => setHost(e.target.value)}
        style={{ width: "100%", marginBottom: 10 }}
      />

      <button onClick={handleConnect}>接続</button>

      {error && <p style={{ color: "red" }}>{error}</p>}
    </div>
  );
};

export default ConnectPage;
