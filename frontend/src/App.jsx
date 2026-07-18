import React, { useState } from 'react';

const API_BASE = 'http://localhost:8000';
const DIMENSIONS = 8;

export default function App() {
  // States for Ingestion
  const [insertInput, setInsertInput] = useState('1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0');
  const [insertStatus, setInsertStatus] = useState({ type: '', msg: '' });

  // States for Searching
  const [searchInput, setSearchInput] = useState('1.1, 2.1, 3.1, 4.1, 5.1, 6.1, 7.1, 8.1');
  const [searchResult, setSearchResult] = useState(null);
  const [searchStatus, setSearchStatus] = useState({ type: '', msg: '' });

  // Helper to parse comma-separated string to float array
  const parseVector = (str) => {
    const arr = str.split(',').map(v => parseFloat(v.trim()));
    if (arr.length !== DIMENSIONS || arr.some(isNaN)) {
      throw new Error(`Must be exactly ${DIMENSIONS} comma-separated numbers.`);
    }
    return arr;
  };

  const handleInsert = async (e) => {
    e.preventDefault();
    setInsertStatus({ type: '', msg: '' });
    try {
      const vector = parseVector(insertInput);
      const res = await fetch(`${API_BASE}/insert`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ vector })
      });
      const data = await res.json();
      if (res.ok && data.success) {
        setInsertStatus({ type: 'success', msg: 'Vector ingested successfully!' });
      } else {
        setInsertStatus({ type: 'error', msg: data.detail || 'Failed to insert vector.' });
      }
    } catch (err) {
      setInsertStatus({ type: 'error', msg: err.message });
    }
  };

  const handleSearch = async (e) => {
    e.preventDefault();
    setSearchStatus({ type: '', msg: '' });
    setSearchResult(null);
    try {
      const vector = parseVector(searchInput);
      const startTime = performance.now();
      
      const res = await fetch(`${API_BASE}/search`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ vector })
      });
      const data = await res.json();
      const endTime = performance.now();

      if (res.ok) {
        setSearchResult({
          nearest: data.nearest_vector,
          distance: data.distance,
          latency: (endTime - startTime).toFixed(2)
        });
      } else {
        setSearchStatus({ type: 'error', msg: data.detail || 'Search failed.' });
      }
    } catch (err) {
      setSearchStatus({ type: 'error', msg: err.message });
    }
  };

  return (
    <div className="min-h-screen bg-slate-900 text-slate-100 font-sans p-8">
      {/* Header */}
      <header className="max-w-5xl mx-auto mb-8 border-b border-slate-800 pb-4">
        <h1 className="text-2xl font-bold tracking-tight text-white">VectorDB Admin Console</h1>
        <p className="text-sm text-slate-400">Minimalist interface for 8-Dimensional Sensor Vector Operations</p>
      </header>

      {/* Main Container */}
      <main className="max-w-5xl mx-auto grid grid-cols-1 md:grid-cols-2 gap-8">
        
        {/* Left Column: Insert */}
        <section className="bg-slate-800 p-6 rounded-lg border border-slate-700 shadow-xl">
          <h2 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <span className="h-2 w-2 rounded-full bg-blue-500"></span>
            Ingest Vector
          </h2>
          <form onSubmit={handleInsert} className="space-y-4">
            <div>
              <label className="block text-xs font-medium text-slate-400 uppercase tracking-wider mb-2">
                Comma-Separated Floats (8 Dimensions)
              </label>
              <input
                type="text"
                className="w-full bg-slate-950 border border-slate-700 rounded p-2.5 text-sm font-mono text-emerald-400 focus:outline-none focus:border-blue-500 transition-colors"
                value={insertInput}
                onChange={(e) => setInsertInput(e.target.value)}
              />
            </div>
            <button
              type="submit"
              className="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium text-sm py-2 px-4 rounded transition-colors shadow"
            >
              Execute Ingest
            </button>
          </form>

          {insertStatus.msg && (
            <div className={`mt-4 p-3 rounded text-sm font-medium ${
              insertStatus.type === 'success' ? 'bg-emerald-950/50 text-emerald-400 border border-emerald-800' : 'bg-rose-950/50 text-rose-400 border border-rose-800'
            }`}>
              {insertStatus.msg}
            </div>
          )}
        </section>

        {/* Right Column: Search */}
        <section className="bg-slate-800 p-6 rounded-lg border border-slate-700 shadow-xl">
          <h2 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <span className="h-2 w-2 rounded-full bg-purple-500"></span>
            Vector Search Sandbox
          </h2>
          <form onSubmit={handleSearch} className="space-y-4">
            <div>
              <label className="block text-xs font-medium text-slate-400 uppercase tracking-wider mb-2">
                Query Vector (8 Dimensions)
              </label>
              <input
                type="text"
                className="w-full bg-slate-950 border border-slate-700 rounded p-2.5 text-sm font-mono text-purple-400 focus:outline-none focus:border-purple-500 transition-colors"
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
              />
            </div>
            <button
              type="submit"
              className="w-full bg-purple-600 hover:bg-purple-700 text-white font-medium text-sm py-2 px-4 rounded transition-colors shadow"
            >
              Run Nearest Neighbor Search
            </button>
          </form>

          {searchStatus.msg && (
            <div className="mt-4 p-3 bg-rose-950/50 text-rose-400 border border-rose-800 rounded text-sm font-medium">
              {searchStatus.msg}
            </div>
          )}

          {searchResult && (
            <div className="mt-6 border-t border-slate-700 pt-4 space-y-3">
              <h3 className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Search Results</h3>
              <div className="bg-slate-950 p-3 rounded border border-slate-800 font-mono text-sm">
                <span className="text-slate-500">Nearest:</span>{' '}
                <span className="text-white">[{searchResult.nearest.map(v => v.toFixed(4)).join(', ')}]</span>
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div className="bg-slate-950 p-3 rounded border border-slate-800">
                  <div className="text-xs text-slate-500 mb-0.5">L2 Distance (Sq)</div>
                  <div className="font-mono text-sm font-bold text-amber-400">{parseFloat(searchResult.distance).toFixed(4)}</div>
                </div>
                <div className="bg-slate-950 p-3 rounded border border-slate-800">
                  <div className="text-xs text-slate-500 mb-0.5">UI Round Trip</div>
                  <div className="font-mono text-sm font-bold text-sky-400">{searchResult.latency} ms</div>
                </div>
              </div>
            </div>
          )}
        </section>

      </main>
    </div>
  );
}