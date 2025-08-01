# GoTunnel Website

A beautiful, animated React/Tailwind website for the GoTunnel project. This website showcases the self-hosted secure tunneling system with interactive animations, code examples, and comprehensive documentation.

## Features

- 🎨 **Modern Design**: Beautiful dark theme with gradient accents
- ✨ **Smooth Animations**: Framer Motion powered animations throughout
- 📱 **Responsive**: Fully responsive design for all devices
- 🎯 **Interactive**: Interactive components and hover effects
- 📝 **Code Examples**: Syntax-highlighted code snippets
- 🚀 **Performance**: Optimized for fast loading and smooth interactions

## Tech Stack

- **React 18** - Modern React with hooks
- **TypeScript** - Type safety and better development experience
- **Tailwind CSS** - Utility-first CSS framework
- **Framer Motion** - Animation library
- **Heroicons** - Beautiful SVG icons
- **React Syntax Highlighter** - Code syntax highlighting
- **React Type Animation** - Typewriter effects

## Getting Started

### Prerequisites

- Node.js 16+ 
- npm or yarn

### Installation

1. Navigate to the website directory:
```bash
cd website
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm start
```

4. Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

### Building for Production

```bash
npm run build
```

This creates a `build` folder with the production-ready files.

## Project Structure

```
website/
├── public/                 # Static assets
├── src/
│   ├── components/         # React components
│   │   ├── Navbar.tsx     # Navigation bar
│   │   ├── Hero.tsx       # Hero section
│   │   ├── Features.tsx   # Features showcase
│   │   ├── HowItWorks.tsx # How it works section
│   │   ├── Architecture.tsx # System architecture
│   │   ├── QuickStart.tsx # Quick start guide
│   │   ├── CodeExamples.tsx # Code examples
│   │   ├── Comparison.tsx # Comparison table
│   │   └── Footer.tsx     # Footer component
│   ├── App.tsx            # Main app component
│   ├── index.tsx          # Entry point
│   └── index.css          # Global styles
├── package.json           # Dependencies and scripts
├── tailwind.config.js     # Tailwind configuration
└── README.md             # This file
```

## Components

### Navbar
Responsive navigation with smooth scrolling and mobile menu.

### Hero
Animated hero section with typewriter effect and floating particles.

### Features
Interactive feature cards with hover animations and descriptions.

### How It Works
Step-by-step explanation with interactive diagrams and terminal examples.

### Architecture
System architecture visualization with layer navigation and technical details.

### Quick Start
Interactive setup guide with terminal examples and configuration files.

### Code Examples
Syntax-highlighted code examples with tabbed navigation.

### Comparison
Feature comparison table with detailed pros/cons analysis.

### Footer
Comprehensive footer with links, social media, and newsletter signup.

## Customization

### Colors
The color scheme is defined in `tailwind.config.js`:

```javascript
colors: {
  primary: { /* Blue gradient */ },
  tunnel: { /* Cyan gradient */ },
  dark: { /* Dark theme colors */ }
}
```

### Animations
Custom animations are defined in `tailwind.config.js` and `index.css`:

```css
@keyframes tunnel-flow {
  0% { left: -100%; }
  100% { left: 100%; }
}
```

### Content
Update the content by modifying the component files in `src/components/`.

## Deployment

### Netlify
1. Connect your repository to Netlify
2. Set build command: `npm run build`
3. Set publish directory: `build`

### Vercel
1. Import your repository to Vercel
2. Vercel will automatically detect the React app
3. Deploy with default settings

### GitHub Pages
1. Add `"homepage": "https://username.github.io/repo-name"` to package.json
2. Install gh-pages: `npm install --save-dev gh-pages`
3. Add deploy script: `"deploy": "gh-pages -d build"`
4. Run: `npm run deploy`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

## Support

For support and questions:
- Create an issue on GitHub
- Join our community discussions
- Check the documentation

---

Built with ❤️ for the GoTunnel community 