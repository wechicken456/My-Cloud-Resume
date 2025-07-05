import { api } from "../api/api";

class VisitorCounter {
    private counterBoard  : HTMLDivElement;
    private visitorCountElement !: HTMLElement;
    constructor() {
        // Visitor Counter
        this.counterBoard = document.createElement("div") as HTMLDivElement;
        this.counterBoard.id = "visitor-counter-board";
        this.counterBoard.innerHTML = `
            <h3>Visitors</h3>
            <span id="visitor-count">Loading...</span>
        `;
    }

    setupVisitorCounter(container: HTMLElement) {
        container.append(this.counterBoard);

        this.visitorCountElement = document.getElementById('visitor-count') as HTMLElement;

        

        // Initialize with current count
        this.updateVisitorCount();

        // Update count periodically
        setInterval(this.updateVisitorCount, 10000);
    }

    private async updateVisitorCount() {
        try {
            const count = await api.fetchVisitorCount();
            if (this.visitorCountElement) {
                this.visitorCountElement.textContent = count;
            }
        } catch (error) {
            console.error('Error updating likes:', error);
            if (this.visitorCountElement) {
                this.visitorCountElement.textContent = 'Error';
            }
        }
    }

    updateVisitorSessionStatus(has_visited: boolean) {
        // If this is a first-time visitor, increment the count
        if (has_visited) {
            api.incrementVisitorCount().then(count => {
                if (this.visitorCountElement) this.visitorCountElement.textContent = count;
            });
        }
    }

}
export var visitorCounter = new VisitorCounter();