import Head from 'next/head'
import { useState } from 'react'

export default function Home() {
  const [activeTab, setActiveTab] = useState('features')

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-purple-900 to-indigo-900">
      <Head>
        <title>GoTunnel - The Open-Source ngrok Killer</title>
        <meta name="description" content="Self-hosted, secure, and blazingly fast tunneling for developers who want control. The open-source alternative to ngrok." />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
        
        {/* Open Graph */}
        <meta property="og:title" content="GoTunnel - The Open-Source ngrok Killer" />
        <meta property="og:description" content="Self-hosted tunneling solution with unlimited tunnels, no quotas, and complete control." />
        <meta property="og:type" content="website" />
        <meta property="og:url" content="https://gotunnel.dev" />
        
        {/* Twitter */}
        <meta name="twitter:card" content="summary_large_image" />
        <meta name="twitter:title" content="GoTunnel - The Open-Source ngrok Killer" />
        <meta name="twitter:description" content="Self-hosted tunneling solution with unlimited tunnels, no quotas, and complete control." />
      </Head>

      {/* Navigation */}
      <nav className="bg-black/20 backdrop-blur-sm border-b border-white/10">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <div className="flex items-center">
              <h1 className="text-2xl font-bold text-white">🚀 GoTunnel</h1>
            </div>
            <div className="hidden md:flex space-x-8">
              <a href="#features" className="text-white/80 hover:text-white transition">Features</a>
              <a href="#comparison" className="text-white/80 hover:text-white transition">vs ngrok</a>
              <a href="#download" className="text-white/80 hover:text-white transition">Download</a>
              <a href="https://github.com/Nom-nom-hub/gotunnel" className="text-white/80 hover:text-white transition">GitHub</a>
            </div>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="py-20 px-4">
        <div className="max-w-7xl mx-auto text-center">
          <h1 className="text-5xl md:text-7xl font-bold text-white mb-6">
            The Open-Source
            <span className="block bg-gradient-to-r from-cyan-400 to-blue-500 bg-clip-text text-transparent">
              ngrok Killer
            </span>
          </h1>
          <p className="text-xl md:text-2xl text-white/80 mb-8 max-w-3xl mx-auto">
            Self-hosted, secure, and blazingly fast tunneling for developers who want control. 
            No limits, no quotas, no data collection.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <a 
              href="https://github.com/Nom-nom-hub/gotunnel/releases"
              className="bg-gradient-to-r from-cyan-500 to-blue-500 text-white px-8 py-4 rounded-lg font-semibold text-lg hover:shadow-lg hover:shadow-cyan-500/25 transition"
            >
              🚀 Download Now
            </a>
            <a 
              href="https://github.com/Nom-nom-hub/gotunnel"
              className="border border-white/20 text-white px-8 py-4 rounded-lg font-semibold text-lg hover:bg-white/10 transition"
            >
              📖 View on GitHub
            </a>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="py-20 px-4">
        <div className="max-w-7xl mx-auto">
          <h2 className="text-4xl font-bold text-white text-center mb-16">Why Choose GoTunnel?</h2>
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
            {[
              {
                icon: "🔒",
                title: "Self-Hosted",
                description: "Complete control over your tunnel infrastructure. Your data, your rules."
              },
              {
                icon: "⚡",
                title: "Blazingly Fast",
                description: "WebSocket-based communication with sub-100ms latency."
              },
              {
                icon: "🔄",
                title: "Unlimited Tunnels",
                description: "No artificial limits or quotas. Create as many tunnels as you need."
              },
              {
                icon: "🖥️",
                title: "GUI Dashboard",
                description: "Beautiful Electron interface with real-time monitoring and logs."
              },
              {
                icon: "🔧",
                title: "Easy Setup",
                description: "One-click installation with automatic dependency management."
              },
              {
                icon: "🌐",
                title: "Cross-Platform",
                description: "Windows, macOS, and Linux support out of the box."
              }
            ].map((feature, index) => (
              <div key={index} className="bg-white/10 backdrop-blur-sm rounded-xl p-6 border border-white/20">
                <div className="text-4xl mb-4">{feature.icon}</div>
                <h3 className="text-xl font-semibold text-white mb-2">{feature.title}</h3>
                <p className="text-white/80">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Comparison Section */}
      <section id="comparison" className="py-20 px-4">
        <div className="max-w-7xl mx-auto">
          <h2 className="text-4xl font-bold text-white text-center mb-16">GoTunnel vs ngrok</h2>
          <div className="bg-white/10 backdrop-blur-sm rounded-xl p-8 border border-white/20">
            <div className="grid md:grid-cols-2 gap-8">
              <div>
                <h3 className="text-2xl font-bold text-white mb-6">GoTunnel ✅</h3>
                <ul className="space-y-3">
                  {[
                    "Self-hosted - Complete control",
                    "Unlimited tunnels - No quotas",
                    "Free forever - No pricing tiers",
                    "Open source - Full transparency",
                    "WebSocket support - Real-time",
                    "GUI dashboard - Easy management",
                    "No data collection - Privacy first",
                    "Custom domains - Full control"
                  ].map((feature, index) => (
                    <li key={index} className="flex items-center text-white/90">
                      <span className="text-green-400 mr-2">✓</span>
                      {feature}
                    </li>
                  ))}
                </ul>
              </div>
              <div>
                <h3 className="text-2xl font-bold text-white mb-6">ngrok ❌</h3>
                <ul className="space-y-3">
                  {[
                    "Cloud-hosted - Vendor lock-in",
                    "Limited tunnels - Paid tiers",
                    "Expensive pricing - $8/month+",
                    "Closed source - No transparency",
                    "HTTP only - Limited protocols",
                    "CLI only - No GUI",
                    "Data collection - Privacy concerns",
                    "Limited domains - Vendor control"
                  ].map((feature, index) => (
                    <li key={index} className="flex items-center text-white/90">
                      <span className="text-red-400 mr-2">✗</span>
                      {feature}
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Download Section */}
      <section id="download" className="py-20 px-4">
        <div className="max-w-7xl mx-auto text-center">
          <h2 className="text-4xl font-bold text-white mb-8">Ready to Get Started?</h2>
          <p className="text-xl text-white/80 mb-12 max-w-2xl mx-auto">
            Download GoTunnel and experience the freedom of self-hosted tunneling. 
            No registration, no limits, no compromises.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <a 
              href="https://github.com/Nom-nom-hub/gotunnel/releases"
              className="bg-gradient-to-r from-green-500 to-emerald-500 text-white px-8 py-4 rounded-lg font-semibold text-lg hover:shadow-lg hover:shadow-green-500/25 transition"
            >
              📦 Download Latest Release
            </a>
            <a 
              href="https://github.com/Nom-nom-hub/gotunnel"
              className="border border-white/20 text-white px-8 py-4 rounded-lg font-semibold text-lg hover:bg-white/10 transition"
            >
              🔧 View Source Code
            </a>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-black/20 border-t border-white/10 py-12">
        <div className="max-w-7xl mx-auto px-4 text-center">
          <p className="text-white/60 mb-4">
            Made with ❤️ by developers, for developers
          </p>
          <p className="text-white/40 text-sm">
            The open-source ngrok killer for developers who want control.
          </p>
        </div>
      </footer>
    </div>
  )
} 