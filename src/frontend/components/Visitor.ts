import { api } from "../api/api";

export function setupVisitorCounter(container: HTMLElement) {
      // Visitor Counter
      const counterBoard = document.createElement("div");
      counterBoard.id = "visitor-counter-board";
      counterBoard.innerHTML = `
          <h3>Visitors</h3>
          <span id="visitor-count">Loading...</span>
      `;
    container.append(counterBoard);

    const visitorCountElement = document.getElementById('visitor-count') as HTMLElement;

    function updateVisitorCount() {
        api.fetchVisitorCount().then(count => {
            if (visitorCountElement) visitorCountElement.textContent = count;
        });
    }
    updateVisitorCount();
    setInterval(updateVisitorCount, 10000);
}
