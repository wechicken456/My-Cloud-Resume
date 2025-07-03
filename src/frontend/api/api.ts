const baseURL = 'https://raf6u1lwte.execute-api.us-east-2.amazonaws.com/test/api';
const visitorCountElement = document.getElementById('visitor-count') as HTMLElement;

export const api = {

    async fetchCount() : Promise<string> {
        const res = await fetch(`${baseURL}/getCount`, {
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

    async fetchLikes() : Promise<string> {
        const res = await fetch(`${baseURL}/getLikes`, {
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

    async incrementLikes() : Promise<string> {
        const res = await fetch(`${baseURL}/incrementLikes`, {
            method: 'POST',
            mode: 'cors',
        })
        .then(response => response.json())
        .catch(error => {
            console.log("Error incrementing likes: ", error);
            return { count: 'Error' };
        });
        return res.count.toString();
    },

    async sendNotification(type: 'like' | 'contact', payload: any) : Promise<void> {
        await fetch(`${baseURL}/notify`, {
            method: 'POST',
            mode: 'cors',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ type, ...payload }),
        });
    },

    async sendContact(form: { name: string; email: string; message: string; recaptcha: string }) : Promise<Response> {
        return fetch(`${baseURL}/contact`, {
            method: 'POST',
            mode: 'cors',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(form),
        });
    }
}
