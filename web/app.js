async function updateStatus() {
    const response = await fetch('/metrics');
    const data = await response.json();
    
    document.getElementById('status-container').innerText = 
        `Active Tunnels: ${data.active_tunnels} | Latency: ${data.latency}ms`;
}

setInterval(updateStatus, 2000);