# Cloud Resume

## Project structure

```
./
|-- public/
│   |- index.html
|-- src/
│   |- data/
│   │   |- resumeData.ts
│   |- main.ts
│   |- style.css
|-- package.json
|-- tsconfig.json
|-- vite.config.ts
|-- tailwind.config.js
```

## Step 1 - HTML, CSS, Typescript

Standard stuffs.

Create app:
```bash
npm create vite@latest cloud-resume
npm install tailwindcss @tailwindcss/vite
npm install typescript
```

Then test the project with
```bash
npm run dev
```

## Step 2 - Create S3 bucket.

Some important configs to look for when creating the bucket:


- Uncheck Block all public access to allow public access (required for static hosting).
- Enable static website hosting from the `Properties` tab of the bucket.

Then upload all the files in the production folder `dist/` to the bucket from the terminal:
```bash
aws s3 cp ./dist s3://tin-cloud-resume/ --recursive
```

## Step 3 - Create Origin Access Identity

By using an OAI, you make sure that viewers can't bypass CloudFront and get the video directly from the S3 bucket. Only the CloudFront OAI can access the file in the S3 bucket. 

1. In the *left navigation plane* of AWS CloudFront console, choose **Origin Access**.
2. Under the **Identities** tab, select **Create origin access identity**.
3. Enter a name for it. 
4. Create

## Step 4 - Create CloudFront Distribution

Some important configs to look for when creating the distribution:

- Origin domain: Select your S3 bucket from the dropdown (e.g., tin-vuong-resume-<random-string>.s3.amazonaws.com). Do **NOT** choose the *website endpoint* option, even if AWS recommends it.
- Leave Origin path blank.
- For **Origin Access**, choose the **Origin access control settings (recommended)**.
    - Select the created OAI from step 3.
    - Copy policy and paste it into the policy of the S3 bucket.
- Redirect HTTP to HTTPS.
- Default cache behavior settings: Leave as default (GET, HEAD allowed).
- Default root object: Enter `index.html`.
- Price class: only North America and Europe.
- Leave Web Application Firewall (WAF) as Do not enable (unless you need it).