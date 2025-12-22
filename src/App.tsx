import React, { useState } from "react";
import ConnectPage from "./pages/ConnectPage.tsx";
import QueryPage from "./pages/QueryPage.tsx";

type ConnectInfo = {
  user: string;
  pass: string;
  host: string;
};

const App: React.FC = () => {
  const [connected, setConnected] = useState(false);
  const [connectInfo, setConnectInfo] = useState<ConnectInfo | null>(null);
  const [dbList, setDbList] = useState<string[]>([]);

  const handleConnected = (info: ConnectInfo, dbs: string[]) => {
    setConnectInfo(info);
    setDbList(dbs);
    setConnected(true);
  };

  return (
    <div style={{ fontFamily: "sans-serif" }}>
      {!connected ? (
        <ConnectPage onConnected={handleConnected} />
      ) : (
        <QueryPage connectInfo={connectInfo!} dbList={dbList} />
      )}
    </div>
  );
};

export default App;
