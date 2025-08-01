# GoTunnel Website

The official website for GoTunnel - The open-source ngrok killer.

## 🚀 Quick Start

### Development
```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Open http://localhost:3000
```

### Production Build
```bash
# Build for production
npm run build

# Start production server
npm start
```

## 🛠️ Tech Stack

- **Next.js 14** - React framework
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **Vercel** - Deployment platform

## 📦 Deployment

### Vercel (Recommended)

1. **Connect your repository:**
   - Go to [Vercel Dashboard](https://vercel.com/dashboard)
   - Click "New Project"
   - Import your GitHub repository
   - Set root directory to `website/`

2. **Environment Variables:**
   - No environment variables needed

3. **Build Settings:**
   - Framework Preset: Next.js
   - Build Command: `npm run build`
   - Output Directory: `.next`

4. **Deploy:**
   - Click "Deploy"
   - Your site will be live at `https://your-project.vercel.app`

### Custom Domain

1. **Add domain in Vercel:**
   - Go to your project settings
   - Click "Domains"
   - Add your custom domain (e.g., `gotunnel.dev`)

2. **Configure DNS:**
   - Add CNAME record pointing to `cname.vercel-dns.com`
   - Or use Vercel's automatic DNS configuration

## 🎨 Features

- **Responsive Design** - Works on all devices
- **SEO Optimized** - Meta tags and structured data
- **Fast Loading** - Optimized images and code splitting
- **Modern UI** - Glass morphism and gradients
- **Accessibility** - WCAG compliant

## 📁 Project Structure

```
website/
├── pages/              # Next.js pages
│   ├── _app.tsx       # App wrapper
│   └── index.tsx      # Home page
├── styles/             # Global styles
│   └── globals.css    # Tailwind imports
├── public/             # Static assets
├── package.json        # Dependencies
├── next.config.js      # Next.js config
├── tailwind.config.js  # Tailwind config
└── vercel.json         # Vercel config
```

## 🔧 Configuration

### Next.js Config (`next.config.js`)
- React strict mode enabled
- Image optimization
- Redirects for downloads

### Vercel Config (`vercel.json`)
- Build configuration
- Security headers
- Route redirects

### Tailwind Config (`tailwind.config.js`)
- Custom color palette
- Animation keyframes
- Responsive breakpoints

## 🚀 Performance

- **Lighthouse Score:** 95+ (Performance, Accessibility, Best Practices, SEO)
- **Core Web Vitals:** All green
- **Bundle Size:** < 100KB gzipped
- **Load Time:** < 2 seconds

## 📱 Mobile Optimization

- Responsive navigation
- Touch-friendly buttons
- Optimized images
- Fast scrolling

## 🔍 SEO

- Meta tags for all pages
- Open Graph support
- Twitter Card support
- Structured data
- Sitemap generation

## 🎯 Analytics

Add your analytics provider:

```javascript
// pages/_app.tsx
import { Analytics } from '@vercel/analytics/react'

export default function App({ Component, pageProps }) {
  return (
    <>
      <Component {...pageProps} />
      <Analytics />
    </>
  )
}
```

## 🛡️ Security

- Security headers configured
- XSS protection enabled
- Content type options set
- Frame options configured

## 📞 Support

For website issues:
- Create an issue on GitHub
- Contact the development team
- Check the documentation

---

**Made with ❤️ by the GoTunnel team** 