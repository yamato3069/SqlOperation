import React, { useState, useMemo } from "react";

type Props = {
  dbList: string[];
  connectInfo: {
    user: string;
    pass: string;
    host: string;
  };
};

type ApiResponse = {
  sql: string;
  data: any[];
};

const QueryPage: React.FC<Props> = ({ dbList }) => {
  const [selectedDB, setSelectedDB] = useState(dbList[0] || "");
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<any[]>([]);
  const [sql, setSql] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  // ページング
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 20;
  const totalPages = Math.ceil(results.length / itemsPerPage);

  const paginatedData = useMemo(() => {
    const start = (currentPage - 1) * itemsPerPage;
    return results.slice(start, start + itemsPerPage);
  }, [results, currentPage]);

  const runQuery = async () => {
    if (!query.trim()) return;

    setLoading(true);
    setError("");
    setResults([]);
    setSql("");

    try {
      const res = await fetch("http://localhost:8080/nl_query", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          query,
          db: selectedDB,
        }),
      });

      if (!res.ok) {
        const msg = await res.text();
        throw new Error(msg);
      }

      const data: ApiResponse = await res.json();
      setSql(data.sql);
      setResults(data.data);
    } catch (err: any) {
      setError(err.message || "エラー");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 900, margin: "30px auto" }}>
      <h2>データベース選択</h2>

      <select
        value={selectedDB}
        onChange={(e) => setSelectedDB(e.target.value)}
        style={{ padding: 8 }}
      >
        {dbList.map((db) => (
          <option key={db} value={db}>
            {db}
          </option>
        ))}
      </select>

      <h2 style={{ marginTop: 30 }}>自然言語クエリ</h2>

      <textarea
        rows={3}
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="例: 山田の注文一覧を表示して"
        style={{ width: "100%", padding: 10 }}
      />

      <button
        onClick={runQuery}
        disabled={loading}
        style={{ marginTop: 10, padding: "10px 20px" }}
      >
        {loading ? "実行中..." : "実行"}
      </button>

      {error && (
        <p style={{ color: "red", marginTop: 10, whiteSpace: "pre-wrap" }}>
          {error}
        </p>
      )}

      {sql && (
        <div style={{ marginTop: 20 }}>
          <strong>生成されたSQL:</strong>
          <pre
            style={{
              background: "#f5f5f5",
              padding: 10,
              whiteSpace: "pre-wrap",
            }}
          >
            {sql}
          </pre>
        </div>
      )}

      {paginatedData.length > 0 && (
        <>
          <table
            style={{
              width: "100%",
              borderCollapse: "collapse",
              marginTop: 20,
            }}
          >
            <thead>
              <tr>
                {Object.keys(paginatedData[0]).map((key) => (
                  <th
                    key={key}
                    style={{
                      borderBottom: "1px solid #ccc",
                      background: "#fafafa",
                      padding: 6,
                    }}
                  >
                    {key}
                  </th>
                ))}
              </tr>
            </thead>

            <tbody>
              {paginatedData.map((row, i) => (
                <tr key={i}>
                  {Object.values(row).map((val, j) => (
                    <td
                      key={j}
                      style={{
                        borderBottom: "1px solid #eee",
                        padding: 6,
                        fontFamily: "monospace",
                      }}
                    >
                      {String(val)}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>

          {/* ページング */}
          <div style={{ marginTop: 15, textAlign: "center" }}>
            <button
              onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
              disabled={currentPage === 1}
            >
              前へ
            </button>
            <span style={{ margin: "0 15px" }}>
              {currentPage} / {totalPages}
            </span>
            <button
              onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))}
              disabled={currentPage === totalPages}
            >
              次へ
            </button>
          </div>
        </>
      )}
    </div>
  );
};

export default QueryPage;
