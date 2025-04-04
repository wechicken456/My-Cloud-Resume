import './style.css';
import { resumeData } from './data';

document.addEventListener('DOMContentLoaded', () => {
    const app = document.getElementById('app') as HTMLElement;
    // Visitor Counter
    const counterBoard = document.createElement("div");
    counterBoard.id = "visitor-counter-board";
    counterBoard.innerHTML = `
        <h3>Visitors</h3>
        <span id="visitor-count">Loading...</span>
    `;
    document.body.append(counterBoard);

    const baseURL = 'https://raf6u1lwte.execute-api.us-east-2.amazonaws.com/test/api';
    const visitorCountElement = document.getElementById('visitor-count') as HTMLElement;

    // check if this client has visited before (per-browser tracking)
    const hasVisited = localStorage.getItem("hasVisited");
    const fetchCount = () => {
        fetch(`${baseURL}/getCount`, {
            method: 'GET',
            mode: 'cors',
        })
        .then(response => response.json())
        .then(data => {
            visitorCountElement.textContent = data.count.toString();
        })
        .catch(error => {
            console.log("Error incrementing count: ", error);
            visitorCountElement.textContent = 'Error...';
        })
    }
    fetchCount();    // fetch the visitor count when the page first loaded
    setInterval(fetchCount, 10000) // then fetch every 10 second to stay up-to-date

    // fetch initial count in the background
    if (!hasVisited) {
        fetch(`${baseURL}/incrementCount`, {
            method: 'POST',
            mode: 'cors',
        })
        .then(response => response.json())
        .then(data => {
            visitorCountElement.textContent = data.count.toString();
            localStorage.setItem('hasVisited', 'true');     // Mark as visisted if it's this browser's first time visting.
        })
        .catch(error => {
            console.log("Error incrementing count: ", error);
            visitorCountElement.textContent = 'Error...';
        })
    }     
    
    // Header
    const header = document.createElement('header');
    header.className = 'header';
    header.innerHTML = `
        <h1>Tin Vuong</h1>
        <div class="contact-info">
            <a href="mailto:${resumeData.contact.email}">${resumeData.contact.email}</a>
            <span>${resumeData.contact.phone}</span>
            <span>${resumeData.contact.location}</span>
            <a href="${resumeData.contact.github}" target="_blank">${resumeData.contact.github}</a>
        </div>
    `;
    app.appendChild(header);

    // Education
    const educationSection = document.createElement('section');
    educationSection.className = 'section';
    educationSection.innerHTML = `
        <h2>Education</h2>
        <div class="entry">
            <h3>${resumeData.education.institution}</h3>
            <span class="date">${resumeData.education.duration}</span>
            <p>${resumeData.education.degree}</p>
            <p style="margin-bottom: 5px;">${resumeData.education.location}</p>
            <ul>
                ${resumeData.education.details.map(detail => `<li>${detail}</li>`).join('')}
            </ul>
        </div>
    `;
    app.appendChild(educationSection);

    // Experience
    const experienceSection = document.createElement('section');
    experienceSection.className = 'section';
    experienceSection.innerHTML = `
        <h2>Work Experience</h2>
        ${resumeData.experience.map(exp => `
            <div class="entry">
                <h3>${exp.title}</h3>
                <p>${exp.organization}</p>
                <span class="date">${exp.duration}</span>
                <p style="margin-bottom: 5px; font-weight: bold;">${exp.location}</p>
                <ul>
                    ${exp.details.map(detail => `<li>${detail}</li>`).join('')}
                </ul>
            </div>
        `).join('')}
    `;
    app.appendChild(experienceSection);

    // Projects
    const projectsSection = document.createElement('section');
    projectsSection.className = 'section';
    projectsSection.innerHTML = `
        <h2>Projects</h2>
        ${resumeData.projects.map(project => `
            <div class="entry">
                <h3>${project.name}</h3>
                <span class="date">${project.date}</span>
                <ul>
                    ${project.details.map(detail => `<li>${detail}</li>`).join('')}
                </ul>
            </div>
        `).join('')}
    `;
    app.appendChild(projectsSection);

    // Activities
    const activitiesSection = document.createElement('section');
    activitiesSection.className = 'section';
    activitiesSection.innerHTML = `
        <h2>Community & Activities</h2>
        <div class="entry">
            <h3>${resumeData.activities.organization}</h3>
            <span class="date">${resumeData.activities.duration}</span>
            <p>${resumeData.activities.roles.join(', ')}</p>
            <ul>
                ${resumeData.activities.details.map(detail => `<li>${detail}</li>`).join('')}
            </ul>
        </div>
    `;
    app.appendChild(activitiesSection);

    // Skills
    const skillsSection = document.createElement('section');
    skillsSection.className = 'section';

    // Group skills by proficiency
    const proficientSkills = resumeData.skills.filter(skill => skill.proficiency === 'proficient');
    const experiencedSkills = resumeData.skills.filter(skill => skill.proficiency === 'experienced');
    const familiarSkills = resumeData.skills.filter(skill => skill.proficiency === 'familiar');

    skillsSection.innerHTML = `
        <div class="skills-header">
            <h2>Skills</h2>
            <div class="skills-legend">
                <span class="legend-item proficient"><span class="legend-color"></span> Proficient</span>
                <span class="legend-item experienced"><span class="legend-color"></span> Experienced</span>
                <span class="legend-item familiar"><span class="legend-color"></span> Familiar</span>
            </div>
        </div>
        <div class="skills-grid">
            <div class="skill-row proficient">
                ${proficientSkills.map(skill => `<span class="skill-tag skill-${skill.proficiency}">${skill.name}</span>`).join('')}
            </div>
            <div class="skill-row experienced">
                ${experiencedSkills.map(skill => `<span class="skill-tag skill-${skill.proficiency}">${skill.name}</span>`).join('')}
            </div>
            <div class="skill-row familiar">
                ${familiarSkills.map(skill => `<span class="skill-tag skill-${skill.proficiency}">${skill.name}</span>`).join('')}
            </div>
        </div>
    `;
    app.appendChild(skillsSection);
});