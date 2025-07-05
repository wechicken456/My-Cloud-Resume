const baseURL = 'https://api.pwnph0fun.com/prod/api';
const visitorCountElement = document.getElementById('visitor-count') as HTMLElement;

export interface SessionStatus {
    has_visited: boolean;
    has_liked: boolean;
}

export const api = {
    async getSession(): Promise<SessionStatus> {
        const res = await fetch(`${baseURL}/session`, {
            method: 'GET',
            mode: 'cors',
            credentials: 'include',
        })
        .then(response => response.json())
        .catch(error => {
            console.log("Error fetching session: ", error);
            throw error;
        });
        return res;
    },

    async fetchVisitorCount() : Promise<string> {
        const res = await fetch(`${baseURL}/getVisitorCount`, {
            method: 'GET',
            mode: 'cors',
        })
        .then(response => response.json())
        .catch(error => {
            console.log("Error fetching count: ", error);
            visitorCountElement.textContent = 'Error...';
        })
        return res.count.toString();
    },

    async incrementVisitorCount() : Promise<string> {
        const res = await fetch(`${baseURL}/incrementVisitorCount`, {
            method: 'POST',
            mode: 'cors',
            credentials: 'include',
        })
        .then(response => response.json())
        .catch(error => {
            console.log("Error incrementing visitor count: ", error);
            return { count: 'Error' };
        });
        return res.count.toString();
    },

    async fetchLikes() : Promise<string> {
        const res = await fetch(`${baseURL}/getLikeCount`, {
            method: 'GET',
            mode: 'cors',
        })
        .then(response => response.json())
        .catch(error => {
            console.log("Error fetching likes: ", error);
            return { count: 'Error' };
        });
        return res.count.toString();
    },

    async toggleLike() : Promise<string> {
        const res = await fetch(`${baseURL}/toggleLike`, {
            method: 'POST',
            mode: 'cors',
            credentials: 'include',
        })
        .then(response => response.json())
        .catch(error => {
            console.log("Error toggling likes: ", error);
            return { count: 'Error' };
        });
        return res.count.toString();
    },

    async sendContact(form: { name: string; email: string; message: string; recaptcha: string }) : Promise<Response> {
        return fetch(`${baseURL}/contact`, {
            method: 'POST',
            mode: 'cors',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(form)
        });
    }
}