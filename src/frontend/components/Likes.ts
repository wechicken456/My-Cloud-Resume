import { api, SessionStatus } from '../api/api';

class LikeCounter {
    private likeDiv : HTMLDivElement;
    private likeCountElement :  HTMLElement;
    private likeBtn : HTMLButtonElement;  
    private likeIcon : SVGElement;
    private isLiked : boolean = false; // Default to not liked
    constructor() {
        this.likeDiv = document.createElement('div') as HTMLDivElement;
        this.likeDiv.id = 'like-counter-board';
        this.likeDiv.innerHTML = `
            <button id="like-btn" aria-label="Like">
                <svg id="like-icon" xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#e74c3c" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
                </svg>
            </button>
            <div id="like-count">Loading...</div>
        `;
    }

    // change the visual state of the like button
    toggleLikeVisual() {
        if (this.isLiked) {
            // Filled red heart
            this.likeIcon.setAttribute('fill', '#e74c3c');
            this.likeIcon.setAttribute('stroke', '#e74c3c');
        } else {
            // Empty heart outline
            this.likeIcon.setAttribute('fill', 'none');
            this.likeIcon.setAttribute('stroke', '#e74c3c');
        }
    }

    setupLikeCounter(container: HTMLElement) {
        container.appendChild(this.likeDiv);

        this.likeCountElement = this.likeDiv.querySelector('#like-count') as HTMLElement;
        this.likeBtn = this.likeDiv.querySelector('#like-btn') as HTMLButtonElement;  
        this.likeIcon = this.likeDiv.querySelector('#like-icon') as SVGElement;
            
        // Initialize visual state based on session
        this.toggleLikeVisual();

        this.likeBtn.onclick = async () => {
            // Toggle like state and update visual
            this.isLiked = !this.isLiked;
            this.toggleLikeVisual();

            const count = await api.toggleLike();
            this.likeCountElement.textContent = count;
        };

        this.updateLikes();
        setInterval(this.updateLikes, 10000);
    }

    // update the like number
    private async updateLikes() {
        try {
            const count = await api.fetchLikes();
            if (this.likeCountElement) {
                this.likeCountElement.textContent = count;
            }
        } catch (error) {
            console.error('Error updating likes:', error);
            if (this.likeCountElement) {
                this.likeCountElement.textContent = 'Error';
            }
        }
    }

    // only called once after the page has loaded
    updateLikeSessionStatus(has_liked: boolean) {
        // If the user has liked, update the visual state   
        if (this.isLiked != has_liked) this.toggleLikeVisual();
        this.isLiked = has_liked;
    }
}

export var likeCounter = new LikeCounter();