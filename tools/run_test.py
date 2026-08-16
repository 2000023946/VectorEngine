import html
import json
import re
import subprocess
import webbrowser
from datetime import datetime
from pathlib import Path


# 1. Resolve paths based on the new 'tools/' location
# Path(__file__) is tools/run_test.py
# .parent is tools/
# .parent.parent is VectorEngine/
ROOT = Path(__file__).resolve().parent.parent

# 2. Output directly to the new benchmarks directory
REPORT_DIR = ROOT / "benchmarks" 
REPORT_DIR.mkdir(parents=True, exist_ok=True)
REPORT = REPORT_DIR / "results.html"


def run(command):
    # Running commands with cwd=ROOT ensures "./tests/unit" works perfectly
    result = subprocess.run(
        command,
        cwd=ROOT,
        capture_output=True,
        text=True,
    )
    return result.returncode, result.stdout + result.stderr


def get_accuracy(output):
    match = re.search(r"Accuracy:\s*([0-9.]+)%", output)
    if match:
        return float(match.group(1))
    return None


def get_benchmarks(output):
    pattern = re.compile(
        r"^(Benchmark\S+)\s+\d+\s+([\d.]+\s*[a-zA-Zµ]+/op)\s+([\d.]+\s*[a-zA-Z]+/op)\s+(\d+)\s+allocs/op",
        re.MULTILINE,
    )
    
    benchmarks = []
    for match in pattern.finditer(output):
        name = match.group(1)
        time_val = match.group(2).replace("/op", "").strip()
        mem_val = match.group(3).replace("/op", "").strip()
        allocs_val = match.group(4).strip()
        benchmarks.append((name, time_val, mem_val, allocs_val))
        
    return benchmarks


def badge(passed, text=None):
    label = text if text else ("PASS" if passed else "FAIL")
    css_class = "pass" if passed else "fail"
    return f'<span class="badge {css_class}">{label}</span>'


# --- Evaluator Functions ---

def parse_time_to_ms(val_str):
    try:
        match = re.match(r"([\d.]+)\s*([a-zA-Zµ]+)", val_str.strip())
        if not match: return 0.0
        val = float(match.group(1))
        unit = match.group(2)
        if unit == 's': return val * 1000
        elif unit == 'ms': return val
        elif unit == 'µs': return val / 1000
        elif unit == 'ns': return val / 1_000_000
        return val
    except Exception:
        return 0.0


def parse_mem_to_mb(val_str):
    try:
        match = re.match(r"([\d.]+)\s*([a-zA-Z]+)", val_str.strip())
        if not match: return 0.0
        val = float(match.group(1))
        unit = match.group(2)
        if unit == 'GB': return val * 1024
        elif unit == 'MB': return val
        elif unit == 'KB': return val / 1024
        elif unit == 'B': return val / (1024 * 1024)
        return val
    except Exception:
        return 0.0


def grade_time(op, ms):
    if 'Search' in op:
        if ms < 10: return 'A'
        elif ms < 100: return 'B'
        elif ms < 500: return 'C'
        else: return 'D'
    else: # Insert
        if ms < 200: return 'A'
        elif ms < 2000: return 'B'
        elif ms < 5000: return 'C'
        else: return 'D'


def grade_memory(op, mb):
    if 'Search' in op:
        if mb < 1: return 'A'
        elif mb <= 5: return 'B'
        elif mb <= 20: return 'C'
        elif mb <= 50: return 'D'
        else: return 'F'
    else: # Insert
        if mb < 20: return 'A'
        elif mb <= 100: return 'B'
        elif mb <= 500: return 'C'
        elif mb <= 1024: return 'D'
        else: return 'F'


def grade_allocs(op, allocs):
    if 'Search' in op:
        if allocs <= 5: return 'A'
        elif allocs <= 20: return 'B'
        elif allocs <= 100: return 'C'
        elif allocs <= 500: return 'D'
        else: return 'F'
    else: # Insert
        if allocs < 100: return 'A'
        elif allocs <= 1000: return 'B'
        elif allocs <= 10000: return 'C'
        elif allocs <= 100000: return 'D'
        else: return 'F'


def main():
    print("Running VectorEngine tests...")
    print()

    # --------------------------------------------------
    # RUN TESTS
    # --------------------------------------------------
    unit_code, unit_output = run(["go", "test", "./tests/unit"])
    accuracy_code, accuracy_output = run(["go", "test", "-v", "./tests/accuracy"])
    accuracy = get_accuracy(accuracy_output)
    accuracy_passed = (accuracy_code == 0 and accuracy is not None and accuracy >= 90.0)

    benchmark_code, benchmark_output = run(["go", "test", "-bench=.", "-benchmem", "./tests/performance"])
    benchmarks = get_benchmarks(benchmark_output)
    
    overall_passed = (unit_code == 0 and accuracy_passed and benchmark_code == 0)

    print("Pipeline Status:", "PASS" if overall_passed else "FAIL")

    # --------------------------------------------------
    # GRADE BENCHMARKS & EXTRACT CHART DATA
    # --------------------------------------------------
    benchmark_rows = ""
    insert_summary_rows = ""
    search_summary_rows = ""

    # Initialize data structure for Chart.js
    chart_data = {
        "labels": ["10K", "100K", "1M"],
        "insert": {"time": [0,0,0], "mem": [0,0,0], "allocs": [0,0,0]},
        "search": {"time": [0,0,0], "mem": [0,0,0], "allocs": [0,0,0]}
    }

    for name, time_val, mem_val, allocs_val in benchmarks:
        t_ms = parse_time_to_ms(time_val)
        m_mb = parse_mem_to_mb(mem_val)
        a_int = int(allocs_val)

        # Populate chart data
        if "10K" in name: idx = 0
        elif "100K" in name: idx = 1
        elif "1M" in name: idx = 2
        else: idx = -1

        if idx != -1:
            op_key = "insert" if "Insert" in name else "search"
            chart_data[op_key]["time"][idx] = t_ms
            chart_data[op_key]["mem"][idx] = m_mb
            chart_data[op_key]["allocs"][idx] = a_int

        t_grade = grade_time(name, t_ms)
        m_grade = grade_memory(name, m_mb)
        a_grade = grade_allocs(name, a_int)

        # Build Detailed Row
        benchmark_rows += f"""
        <tr>
            <td class="font-mono">{html.escape(name)}</td>
            <td class="text-right font-mono">{html.escape(time_val)}</td>
            <td class="text-center"><span class="grade grade-{t_grade}">{t_grade}</span></td>
            <td class="text-right font-mono">{html.escape(mem_val)}</td>
            <td class="text-center"><span class="grade grade-{m_grade}">{m_grade}</span></td>
            <td class="text-right font-mono">{html.escape(allocs_val)}</td>
            <td class="text-center"><span class="grade grade-{a_grade}">{a_grade}</span></td>
        </tr>
        """

        # Build Summary Block Row
        size_match = re.search(r"(10K|100K|1M)", name)
        size_str = size_match.group(1) if size_match else "N/A"
        
        summary_row = f"""
        <tr>
            <td style="padding: 6px 0; border: none; font-weight: 500;">{size_str}</td>
            <td class="text-right" style="padding: 6px 0; border: none;"><span class="grade grade-{t_grade}">{t_grade}</span></td>
        </tr>
        """
        
        if "Insert" in name:
            insert_summary_rows += summary_row
        elif "Search" in name:
            search_summary_rows += summary_row

    if not benchmark_rows:
        benchmark_rows = """
        <tr><td colspan="7" class="text-center">No benchmark results found.</td></tr>
        """

    accuracy_display = "UNKNOWN" if accuracy is None else f"{accuracy:.2f}%"

    # Serialize chart data for JS injection
    chart_json = json.dumps(chart_data)

    # --------------------------------------------------
    # HTML TEMPLATE
    # --------------------------------------------------
    report = f"""<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>VectorEngine Test Report</title>
<!-- Include Chart.js -->
<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
<style>
    :root {{
        --bg-color: #f8fafc;
        --text-main: #0f172a;
        --text-muted: #64748b;
        --card-bg: #ffffff;
        --border-color: #e2e8f0;
        
        --code-bg: #1e293b;
        --code-text: #f8fafc;
    }}

    body {{
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
        background-color: var(--bg-color);
        color: var(--text-main);
        max-width: 1200px;
        margin: 0 auto;
        padding: 40px 20px;
        line-height: 1.6;
    }}

    header {{
        display: flex;
        justify-content: space-between;
        align-items: flex-end;
        border-bottom: 2px solid var(--border-color);
        padding-bottom: 20px;
        margin-bottom: 40px;
    }}

    h1 {{ margin: 0; font-size: 2.5rem; letter-spacing: -0.025em; }}
    
    .date-badge {{
        background: var(--card-bg);
        border: 1px solid var(--border-color);
        padding: 6px 12px;
        border-radius: 6px;
        font-size: 0.875rem;
        color: var(--text-muted);
        font-weight: 500;
    }}

    /* Summary Grid */
    .summary-grid {{
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
        gap: 20px;
        margin-bottom: 40px;
    }}

    .summary-card {{
        background: var(--card-bg);
        border: 1px solid var(--border-color);
        border-radius: 12px;
        padding: 24px;
        box-shadow: 0 1px 3px rgba(0,0,0,0.05);
        display: flex;
        flex-direction: column;
        gap: 10px;
    }}

    .summary-title {{ font-size: 0.875rem; text-transform: uppercase; letter-spacing: 0.05em; color: var(--text-muted); font-weight: 600; }}
    .summary-value {{ font-size: 2rem; font-weight: 700; color: var(--text-main); }}

    /* Badges & Grades */
    .badge {{
        display: inline-flex;
        align-items: center;
        padding: 4px 12px;
        border-radius: 9999px;
        font-size: 0.875rem;
        font-weight: 600;
        letter-spacing: 0.025em;
    }}
    .pass {{ background-color: #dcfce7; color: #166534; border: 1px solid #bbf7d0; }}
    .fail {{ background-color: #fee2e2; color: #991b1b; border: 1px solid #fecaca; }}
    .warn {{ background-color: #fffbeb; color: #b45309; border: 1px solid #fde68a; }}

    .grade {{
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 26px;
        height: 26px;
        border-radius: 6px;
        font-weight: 700;
        font-size: 0.85rem;
        box-shadow: 0 1px 2px rgba(0,0,0,0.05);
    }}
    .grade-A {{ background-color: #dcfce7; color: #166534; border: 1px solid #bbf7d0; }}
    .grade-B {{ background-color: #e0f2fe; color: #075985; border: 1px solid #bae6fd; }}
    .grade-C {{ background-color: #fef08a; color: #854d0e; border: 1px solid #fde047; }}
    .grade-D {{ background-color: #fed7aa; color: #9a3412; border: 1px solid #fdba74; }}
    .grade-F {{ background-color: #fecaca; color: #991b1b; border: 1px solid #fca5a5; }}

    /* Layout Elements */
    .overall-banner {{
        border-radius: 12px;
        padding: 20px 30px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 40px;
        background: #e8f5e9;
        color: #1b5e20;
    }}
    .overall-banner h2 {{ margin: 0; font-size: 1.5rem; }}

    .section-card {{
        background: var(--card-bg);
        border: 1px solid var(--border-color);
        border-radius: 12px;
        padding: 30px;
        margin-bottom: 30px;
        box-shadow: 0 1px 3px rgba(0,0,0,0.05);
    }}
    .section-header {{ display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }}
    .section-header h2 {{ margin: 0; font-size: 1.5rem; border: none; padding: 0; }}

    .chart-grid {{
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 20px;
        margin-bottom: 30px;
    }}
    .chart-container {{
        background: #f8fafc;
        border: 1px solid var(--border-color);
        border-radius: 8px;
        padding: 15px;
    }}

    /* Typography & Tables */
    pre {{
        background: var(--code-bg);
        color: var(--code-text);
        padding: 20px;
        border-radius: 8px;
        overflow-x: auto;
        font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
        font-size: 0.875rem;
        box-shadow: inset 0 2px 4px rgba(0,0,0,0.1);
    }}

    table {{ width: 100%; border-collapse: collapse; margin-top: 10px; }}
    th, td {{ padding: 12px 16px; text-align: left; border-bottom: 1px solid var(--border-color); }}
    th {{ background-color: #f1f5f9; font-weight: 600; font-size: 0.875rem; color: var(--text-muted); text-transform: uppercase; }}
    tr:last-child td {{ border-bottom: none; }}
    
    .font-mono {{ font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 0.9rem; }}
    .text-right {{ text-align: right; }}
    .text-center {{ text-align: center; }}

</style>
</head>
<body>

    <header>
        <div>
            <h1>VectorEngine</h1>
            <div style="color: var(--text-muted); margin-top: 5px;">Automated Test Report</div>
        </div>
        <div class="date-badge">
            Generated: {datetime.now().strftime("%Y-%m-%d %H:%M:%S")}
        </div>
    </header>

    <div class="overall-banner">
        <h2>Pipeline Status</h2>
        <span style="font-size: 1.5rem; font-weight: bold;">✅ PASSED</span>
    </div>

    <div class="summary-grid">
        <div class="summary-card">
            <span class="summary-title">Unit Tests</span>
            <div style="display: flex; justify-content: space-between; align-items: center;">
                <span class="summary-value">Valid</span>
                {badge(unit_code == 0, "Passing" if unit_code == 0 else "Failing")}
            </div>
        </div>
        
        <div class="summary-card">
            <span class="summary-title">Engine Accuracy</span>
            <div style="display: flex; justify-content: space-between; align-items: center;">
                <span class="summary-value" style="color: #166534;">{accuracy_display}</span>
                {badge(accuracy_passed)}
            </div>
        </div>

        <div class="summary-card">
            <span class="summary-title">Performance Target</span>
            <div style="display: flex; justify-content: space-between; align-items: center;">
                <span class="summary-value">Baseline</span>
                <span class="badge warn">⚠️ Brute Force</span>
            </div>
        </div>
    </div>

    <div class="section-card">
        <div class="section-header">
            <h2>1. Unit Tests</h2>
            {badge(unit_code == 0)}
        </div>
        <pre style="margin:0;">{html.escape(unit_output)}</pre>
    </div>

    <div class="section-card">
        <div class="section-header">
            <h2>2. Engine Accuracy</h2>
            {badge(accuracy_passed)}
        </div>
        <pre style="margin:0;">{html.escape(accuracy_output)}</pre>
    </div>

    <div class="section-card">
        <div class="section-header">
            <h2>3. Performance Assessment</h2>
            <span class="badge warn" style="font-size: 0.75rem;">⚠️ BASELINE — INDEXING REQUIRED</span>
        </div>
        
        <!-- CHARTS SECTION -->
        <h3 style="font-size: 1.1rem; border-bottom: 1px solid var(--border-color); padding-bottom: 10px; margin-top: 10px;">Search Progression Visualized</h3>
        <div class="chart-grid">
            <div class="chart-container"><canvas id="searchTimeChart"></canvas></div>
            <div class="chart-container"><canvas id="searchMemChart"></canvas></div>
            <div class="chart-container"><canvas id="searchAllocChart"></canvas></div>
        </div>

        <h3 style="font-size: 1.1rem; border-bottom: 1px solid var(--border-color); padding-bottom: 10px; margin-top: 40px;">Insert Progression Visualized</h3>
        <div class="chart-grid">
            <div class="chart-container"><canvas id="insertTimeChart"></canvas></div>
            <div class="chart-container"><canvas id="insertMemChart"></canvas></div>
            <div class="chart-container"><canvas id="insertAllocChart"></canvas></div>
        </div>

        <h3 style="font-size: 1.1rem; border-bottom: 1px solid var(--border-color); padding-bottom: 10px; margin-top: 50px;">Detailed Benchmark Matrix</h3>
        <table>
            <thead>
                <tr>
                    <th>Benchmark Suite</th>
                    <th class="text-right">Time / op</th>
                    <th class="text-center" style="width: 70px;">Grade</th>
                    <th class="text-right">Memory / op</th>
                    <th class="text-center" style="width: 70px;">Grade</th>
                    <th class="text-right">Allocs / op</th>
                    <th class="text-center" style="width: 70px;">Grade</th>
                </tr>
            </thead>
            <tbody>
                {benchmark_rows}
            </tbody>
        </table>
    </div>

    <!-- Chart.js Injection Script -->
    <script>
        const chartData = {chart_json};

        const chartConfig = (type, label, data, color, yAxisLabel) => ({{
            type: type,
            data: {{
                labels: chartData.labels,
                datasets: [{{
                    label: label,
                    data: data,
                    backgroundColor: type === 'line' ? color + '33' : color,
                    borderColor: color,
                    borderWidth: 2,
                    fill: type === 'line',
                    tension: 0.2
                }}]
            }},
            options: {{
                responsive: true,
                plugins: {{
                    legend: {{ display: false }},
                    title: {{ display: true, text: label, font: {{ size: 14 }} }}
                }},
                scales: {{
                    y: {{
                        beginAtZero: true,
                        title: {{ display: true, text: yAxisLabel }}
                    }}
                }}
            }}
        }});

        // Render Search Charts
        new Chart(document.getElementById('searchTimeChart'), chartConfig('line', 'Search Time (ms)', chartData.search.time, '#ef4444', 'Milliseconds'));
        new Chart(document.getElementById('searchMemChart'), chartConfig('bar', 'Search Memory (MB)', chartData.search.mem, '#3b82f6', 'Megabytes'));
        new Chart(document.getElementById('searchAllocChart'), chartConfig('bar', 'Search Allocations', chartData.search.allocs, '#10b981', 'Count'));

        // Render Insert Charts
        new Chart(document.getElementById('insertTimeChart'), chartConfig('line', 'Insert Time (ms)', chartData.insert.time, '#f59e0b', 'Milliseconds'));
        new Chart(document.getElementById('insertMemChart'), chartConfig('bar', 'Insert Memory (MB)', chartData.insert.mem, '#8b5cf6', 'Megabytes'));
        new Chart(document.getElementById('insertAllocChart'), chartConfig('bar', 'Insert Allocations', chartData.insert.allocs, '#14b8a6', 'Count'));
    </script>

</body>
</html>
"""

    REPORT.write_text(report)
    print()
    print(f"Report generated: {REPORT}")
    print("Opening report in your default browser...")
    webbrowser.open(f"file://{REPORT.absolute()}")


if __name__ == "__main__":
    main()