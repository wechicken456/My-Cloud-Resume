import { api } from '../api/api';

export function setupLikesCounter(container: HTMLElement) {
    const likesDiv = document.createElement('div');
    likesDiv.id = 'likes-counter-board';
    likesDiv.innerHTML = `
        <button id="like-btn" aria-label="Like">
            <svg id="like-icon" xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#e74c3c" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
            </svg>
        </button>
        <div id="likes-count">0</div>
    `;
    container.appendChild(likesDiv);

    const likesCount = likesDiv.querySelector('#likes-count') as HTMLElement;
    const likeBtn = likesDiv.querySelector('#like-btn') as HTMLButtonElement;  
    const likeIcon = likesDiv.querySelector('#like-icon') as SVGElement;
        
    let isLiked = false;

    function toggleLikeVisual() {
        if (isLiked) {
            // Filled red heart
            likeIcon.setAttribute('fill', '#e74c3c');
            likeIcon.setAttribute('stroke', '#e74c3c');
        } else {
            // Empty heart outline
            likeIcon.setAttribute('fill', 'none');
            likeIcon.setAttribute('stroke', '#e74c3c');
        }
    }

    likeBtn.onclick = async () => {
        // Toggle like state and update visual
        isLiked = !isLiked;
        toggleLikeVisual();

        const count = await api.toggleLike();
        likesCount.textContent = count;
    };

    function updateLikes() {
        api.fetchLikes().then(count => likesCount.textContent = count);
    }

    updateLikes();
    setInterval(updateLikes, 10000);
}