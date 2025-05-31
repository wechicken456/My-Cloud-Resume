export interface ResumeData {
    contact: {
        email: string;
        phone: string;
        location: string;
        github: string;
    };
    education: {
        institution: string;
        duration: string;
        degree: string;
        location: string;
        details: string[];
    };
    experience: Array<{
        title: string;
        organization: string;
        duration: string;
        location: string;
        details: string[];
    }>;
    projects: Array<{
        name: string;
        date: string;
        details: string[];
    }>;
    activities: {
        organization: string;
        duration: string;
        roles: string[];
        details: string[];
    };
    skills: Array<{ name: string, proficiency: string }>;
}

export const resumeData: ResumeData = {
    contact: {
        email: "tinvuong2003@gmail.com",
        phone: "513-914-8401",
        location: "Athens, OH",
        github: "https://github.com/wechicken456"
    },
    education: {
        institution: "Ohio University",
        duration: "August 2022 - Present",
        degree: "Bachelor of Science in Advanced Computing (Honors Tutorial College)",
        location: "Athens, OH",
        details: [
            "As a student in the Honors Tutorial College, I benefit from a distinctive academic experience: personalized, one-on-one tutorials with professors, allowing me to explore subjects of my choosing in depth.",
            "My past tutorials focused on Post-Quantum Cryptography and Linux Kernel Security.",
            "Current tutorial focuses on developing a physics-informed machine learning model using the transformer architecture in PyTorch for aviation engine failure prediction.",
            "Notable coursework: Database Systems, Parallel Computing, Operating Systems, Software Security.",
            "GPA: 3.936, expected graduation May 2026",
            "Dean's List Spring 2024; OHIO Premier Scholarship (full-ride)"
        ]
    },
    experience: [
        {
            title: "Research Assistant",
            organization: "Ohio University, Russ College of Engineering",
            duration: "May 2023 - Present",
            location: "Athens, OH",
            details: [
                "Studied and taught Fully Homomorphic Encryption (FHE) to team members to enhance understanding and application of cryptographic principles in multi-party computation.",
                "Tested and verified our hardware monitors developed in Verilog to ensure system reliability and security."
            ]
        },
        {
            title: "Software Developer",
            organization: "Ohio University, Voinovich School",
            duration: "May 2023 - Present",
            location: "Athens, OH",
            details: [
                "Collaborated with the Director of GVS/ITS and other Ohio University students to develop a AI-powered image search system using Huggingface APIs.",
                "Developed a Retrieval-Augmented Generation (RAG) chatbot integrating OpenAI embeddings, Supabase (PostgreSQL), and PydanticAI to enable querying specialized knowledge bases built from crawled website content.",
                "Implemented an automated pipeline to crawl web pages, extract and embed content, store embeddings with metadata, and retrieve relevant information using agentic LLM processing via PydanticAI."
            ]
        }
    ],
    projects: [
        {
            name: "Microservices-based Banking App in Golang",
            date: "April 2025 - Present",
            details: [
                "Architected a Golang-based Banking Application utilizing a microservices architecture with the Controller-Service-Repository design pattern, facilitating inter-service communication with gRPC, and containerizing microservices and databases with Docker.",
                "Engineered PostgreSQL database solutions for backend services, ensuring ACID transactional integrity for financial operations, managing migrations with goose, and generating type-safe SQL queries with SQLC.",
                "Developed secure JWT-based authentication at the JSON-based RESTful API Gateway, incorporating refresh tokens and fingerprinting adhering to OWASP principles.",
            ]
        },
        {
            name: "Aviation Engine Failure Detection Model",
            date: "March 2025 - Current",
            details: [
                "Developed a LSTM model in PyTorch for aviation engine failure prediction.",
                "Automated synthetic flights using the XPlane SDK to record and process data for model training.",
            ]
        },
        {
            name: "AWS Cloud Resume",
            date: "March 2025",
            details: [
                "Designed and deployed a cloud-hosted resume website using AWS services, including S3 for static hosting, CloudFront for CDN distribution, Route 53 for custom domain management (pwnph0fun.com), and ACM for SSL certification.",
                "Developed a serverless visitor counter with DynamoDB, Go-based Lambda, and API Gateway, integrated with a TypeScript frontend, featuring automated CI/CD pipelines via GitHub Actions and unit testing with Go mocks."
            ]
        },
        {
            name: "Agentic RAG Chatbot",
            date: "Feb 2025",
            details: [
                "Developed a Retrieval-Augmented Generation (RAG) chatbot integrating OpenAI embeddings, Supabase (PostgreSQL), and PydanticAI to enable querying specialized knowledge bases built from crawled website content.",
                "Implemented an automated pipeline to crawl web pages, extract and embed content, store embeddings with metadata, and retrieve relevant information using agentic LLM processing via PydanticAI."
            ]
        },
        {
            name: "Kubernetes The Hard Way",
            date: "Feb 2025",
            details: [
                "Provisioned and configured a Kubernetes cluster on VirtualBox, deploying core components including etcd, kube-apiserver, kube-controller-manager, kube-scheduler, kubelet, and kube-proxy.",
                "Validated Kubernetes functionality through pod deployment, service discovery, and network policy implementation gaining experience in cluster operations and troubleshooting."
            ]
        },
        {
            name: "File Storage Server in Golang",
            date: "Feb 2025",
            details: [
                "Developed a file storage server enabling media uploads to AWS S3 buckets, integrating FastStart processing for video range requests and implementing secure user authentication using JWT and SQLite3.",
                "Configured AWS S3 policies and AWS CloudFront for scalable and resilient content delivery, leveraging Go’s AWS SDK for interaction with AWS services."
            ]
        },
        {
            name: "TinyGator: A Go-Based Blog Aggregator",
            date: "Jan 2025",
            details: [
                "Built CLI app in Go to aggregate RSS feeds with PostgreSQLBuilt a CLI application in Go to aggregate, browse, and manage RSS feeds with PostgreSQL-backed storage.",
                "Leveraged sqlc for type-safe SQL query generation and designed a customizable configuration system using JSON.",
                "Implemented user authentication, user management, and efficient feed parsing to enhance usability and scalability."
            ]
        },
        {
            name: "MapReduce Implementation in Golang",
            date: "October 2024",
            details: [
                "Implemented a relaxed MapReduce framework in Golang, leveraging interfaces, synchronization primitives, goroutines, and channels for efficient task coordination and RFC over Unix sockets for communication.",
                "Currently implementing the Raft consensus algorithm for fault-tolerant coordination."
            ]
        },
        {
            name: "pwn.college Hacking Challenges",
            date: "October 2024",
            details: ["Completed the Sandboxing challenges of pwn.college, covering seccomp and Linux namespaces and how to escape from them. Completed 17/22 Race Condition challenges, covering how to exploit Linux race conditions."]
        },
        {
            name: "Kaggle House Prices Prediction",
            date: "August 2024",
            details: [
                "Built a MLP regression neural network in PyTorch to predict Kaggle house prices, ranking 1221/4389.",
                "Used Pandas for data processing, applied k-fold validation, and integrated batch normalization to improve model performance and mitigate overfitting."
            ]
        },
        {
            name: "Security Exploit Challenge",
            date: "December 2023",
            details: [
                "Developed a heap overflow Capture The Flag (CTF) challenge in C, inspired by the CVE-2021-4034 local privilege escalation vulnerability on Linux.",
                "Published a solution walkthrough on Youtube and a detailed writeup on Github."
            ]
        }
    ],
    activities: {
        organization: "Association for Computing Machinery",
        duration: "February 2023 - Present",
        roles: ["Director of Competitive Programming", "Treasurer", "Member"],
        details: [
            "1st in 2022 UC Hackathon CTF.",
            "21st in 2023 ICPC Mid-Atlantic Regional.",
            "33rd/648 globally (9th undergrad) in 2024 BuckeyeCTF.",
            "30th/688 globally (12th undergrad) in 2022 BuckeyeCTF."
        ]
    },
    skills: [
        { name: "gRPC", proficiency: "proficient" },
        { name: "Huggingface", proficiency: "experienced" },
        { name: "Ollama", proficiency: "experienced" },
        { name: "Supabase", proficiency: "experienced" },
        { name: "n8n", proficiency: "experienced" },
        { name: "Github Actions CI/CD", proficiency: "experienced" },
        { name: "Golang", proficiency: "proficient" },
        { name: "C", proficiency: "proficient" },
        { name: "C++", proficiency: "experienced" },
        { name: "Javascript", proficiency: "experienced" },
        { name: "Python", proficiency: "proficient" },
        { name: "Docker", proficiency: "proficient" },
        { name: "Kubernetes", proficiency: "experienced" },
        { name: "AWS S3", proficiency: "experienced" },
        { name: "AWS CloudFront", proficiency: "experienced" },
        { name: "SQLite3", proficiency: "experienced" },
        { name: "PostgreSQL", proficiency: "proficient" },
        { name: "sqlc", proficiency: "experienced" },
        { name: "Express", proficiency: "familiar" },
        { name: "React", proficiency: "experienced" },
        { name: "Pytorch", proficiency: "experienced" },
        { name: "Git", proficiency: "proficient" },
        { name: "Ghidra", proficiency: "proficient" },
        { name: "Bash", proficiency: "proficient" },
        { name: "Goose Migration", proficiency: "familiar" },
        { name: "Linux Kernel", proficiency: "experienced" },
        { name: "OpenMPI", proficiency: "familiar" },
        { name: "Algorithms", proficiency: "proficient" },
        { name: "Cryptography", proficiency: "proficient" },
        { name: "TypeScript", proficiency: "experienced" },
        { name: "Vite", proficiency: "experienced" },
        { name: "gRPC", proficiency: "experienced" },
        { name: "x86 Assembly", proficiency: "experienced" },
    ]
};
