import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { 
  ServerIcon,
  ComputerDesktopIcon,
  GlobeAltIcon,
  LockClosedIcon,
  CogIcon,
  ChartBarIcon
} from '@heroicons/react/24/outline';

const Architecture: React.FC = () => {
  const [activeLayer, setActiveLayer] = useState(0);

  const layers = [
    {
      name: "Client Layer",
      components: [
        { name: "GoTunnel Client", description: "CLI tool for establishing tunnels", icon: ComputerDesktopIcon },
        { name: "Local Service", description: "Your application (port 3000)", icon: CogIcon },
        { name: "Authentication", description: "Token-based auth", icon: LockClosedIcon }
      ],
      color: "from-blue-500 to-cyan-600"
    },
    {
      name: "Network Layer",
      components: [
        { name: "WebSocket Connection", description: "Bidirectional communication", icon: GlobeAltIcon },
        { name: "TLS Encryption", description: "End-to-end security", icon: LockClosedIcon },
        { name: "HTTP/TCP Proxy", description: "Traffic forwarding", icon: CogIcon }
      ],
      color: "from-green-500 to-emerald-600"
    },
    {
      name: "Server Layer",
      components: [
        { name: "GoTunnel Server", description: "Main server application", icon: ServerIcon },
        { name: "Tunnel Manager", description: "Connection tracking", icon: ChartBarIcon },
        { name: "Subdomain Router", description: "Request routing", icon: CogIcon }
      ],
      color: "from-purple-500 to-pink-600"
    }
  ];

  return (
    <section id="architecture" className="py-20 bg-dark-800/50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
          className="text-center mb-16"
        >
          <h2 className="text-4xl md:text-5xl font-bold mb-6">
            System <span className="text-gradient">Architecture</span>
          </h2>
          <p className="text-xl text-gray-300 max-w-3xl mx-auto">
            A modular, scalable architecture designed for security, performance, and ease of deployment.
          </p>
        </motion.div>

        {/* Layer Navigation */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.2 }}
          viewport={{ once: true }}
          className="flex justify-center mb-12"
        >
          <div className="flex space-x-2 bg-dark-700 rounded-lg p-1">
            {layers.map((layer, index) => (
              <motion.button
                key={index}
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                onClick={() => setActiveLayer(index)}
                className={`px-4 py-2 rounded-md text-sm font-medium transition-all duration-200 ${
                  activeLayer === index
                    ? 'bg-gradient-to-r from-tunnel-500 to-primary-600 text-white'
                    : 'text-gray-400 hover:text-white'
                }`}
              >
                {layer.name}
              </motion.button>
            ))}
          </div>
        </motion.div>

        {/* Architecture Diagram */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-start">
          {/* Layer Details */}
          <motion.div
            initial={{ opacity: 0, x: -30 }}
            whileInView={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
            className="space-y-6"
          >
            <div className={`p-6 rounded-xl border bg-gradient-to-r ${layers[activeLayer].color}/10 border-${layers[activeLayer].color.split('-')[1]}-500/30`}>
              <h3 className="text-xl font-semibold mb-4 text-white">
                {layers[activeLayer].name}
              </h3>
              <div className="space-y-4">
                {layers[activeLayer].components.map((component, index) => (
                  <motion.div
                    key={index}
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ duration: 0.5, delay: index * 0.1 }}
                    className="flex items-center space-x-3"
                  >
                    <div className={`w-10 h-10 rounded-lg bg-gradient-to-r ${layers[activeLayer].color} flex items-center justify-center flex-shrink-0`}>
                      <component.icon className="w-5 h-5 text-white" />
                    </div>
                    <div>
                      <h4 className="text-white font-medium">{component.name}</h4>
                      <p className="text-gray-400 text-sm">{component.description}</p>
                    </div>
                  </motion.div>
                ))}
              </div>
            </div>
          </motion.div>

          {/* Visual Architecture */}
          <motion.div
            initial={{ opacity: 0, x: 30 }}
            whileInView={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
            className="relative"
          >
            <div className="bg-dark-800 rounded-2xl p-8 border border-gray-700">
              {/* Client Layer */}
              <motion.div
                animate={{
                  scale: activeLayer === 0 ? 1.05 : 1,
                  boxShadow: activeLayer === 0 ? "0 0 20px rgba(14, 165, 233, 0.3)" : "none"
                }}
                transition={{ duration: 0.3 }}
                className="mb-6"
              >
                <div className="flex items-center space-x-3 mb-4">
                  <div className="w-8 h-8 bg-gradient-to-r from-blue-500 to-cyan-600 rounded-lg flex items-center justify-center">
                    <ComputerDesktopIcon className="w-4 h-4 text-white" />
                  </div>
                  <h4 className="text-white font-semibold">Client Layer</h4>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="bg-blue-500/20 border border-blue-500/30 rounded-lg p-3">
                    <h5 className="text-blue-300 text-sm font-medium">GoTunnel Client</h5>
                    <p className="text-gray-400 text-xs">CLI tool</p>
                  </div>
                  <div className="bg-blue-500/20 border border-blue-500/30 rounded-lg p-3">
                    <h5 className="text-blue-300 text-sm font-medium">Local Service</h5>
                    <p className="text-gray-400 text-xs">Port 3000</p>
                  </div>
                </div>
              </motion.div>

              {/* Connection */}
              <div className="relative mb-6">
                <div className="h-2 bg-gradient-to-r from-blue-500 to-green-500 rounded-full"></div>
                <motion.div
                  animate={{
                    scaleX: activeLayer >= 1 ? [0, 1, 0] : 0
                  }}
                  transition={{ duration: 2, repeat: activeLayer >= 1 ? Infinity : 0 }}
                  className="absolute top-0 left-0 h-2 bg-tunnel-400 rounded-full origin-left"
                  style={{ width: '100%' }}
                />
              </div>

              {/* Network Layer */}
              <motion.div
                animate={{
                  scale: activeLayer === 1 ? 1.05 : 1,
                  boxShadow: activeLayer === 1 ? "0 0 20px rgba(34, 197, 94, 0.3)" : "none"
                }}
                transition={{ duration: 0.3 }}
                className="mb-6"
              >
                <div className="flex items-center space-x-3 mb-4">
                  <div className="w-8 h-8 bg-gradient-to-r from-green-500 to-emerald-600 rounded-lg flex items-center justify-center">
                    <GlobeAltIcon className="w-4 h-4 text-white" />
                  </div>
                  <h4 className="text-white font-semibold">Network Layer</h4>
                </div>
                <div className="grid grid-cols-3 gap-3">
                  <div className="bg-green-500/20 border border-green-500/30 rounded-lg p-3">
                    <h5 className="text-green-300 text-sm font-medium">WebSocket</h5>
                    <p className="text-gray-400 text-xs">Bidirectional</p>
                  </div>
                  <div className="bg-green-500/20 border border-green-500/30 rounded-lg p-3">
                    <h5 className="text-green-300 text-sm font-medium">TLS</h5>
                    <p className="text-gray-400 text-xs">Encryption</p>
                  </div>
                  <div className="bg-green-500/20 border border-green-500/30 rounded-lg p-3">
                    <h5 className="text-green-300 text-sm font-medium">Proxy</h5>
                    <p className="text-gray-400 text-xs">HTTP/TCP</p>
                  </div>
                </div>
              </motion.div>

              {/* Connection */}
              <div className="relative mb-6">
                <div className="h-2 bg-gradient-to-r from-green-500 to-purple-500 rounded-full"></div>
                <motion.div
                  animate={{
                    scaleX: activeLayer >= 2 ? [0, 1, 0] : 0
                  }}
                  transition={{ duration: 2, repeat: activeLayer >= 2 ? Infinity : 0 }}
                  className="absolute top-0 left-0 h-2 bg-tunnel-400 rounded-full origin-left"
                  style={{ width: '100%' }}
                />
              </div>

              {/* Server Layer */}
              <motion.div
                animate={{
                  scale: activeLayer === 2 ? 1.05 : 1,
                  boxShadow: activeLayer === 2 ? "0 0 20px rgba(168, 85, 247, 0.3)" : "none"
                }}
                transition={{ duration: 0.3 }}
              >
                <div className="flex items-center space-x-3 mb-4">
                  <div className="w-8 h-8 bg-gradient-to-r from-purple-500 to-pink-600 rounded-lg flex items-center justify-center">
                    <ServerIcon className="w-4 h-4 text-white" />
                  </div>
                  <h4 className="text-white font-semibold">Server Layer</h4>
                </div>
                <div className="grid grid-cols-3 gap-3">
                  <div className="bg-purple-500/20 border border-purple-500/30 rounded-lg p-3">
                    <h5 className="text-purple-300 text-sm font-medium">Server</h5>
                    <p className="text-gray-400 text-xs">Main app</p>
                  </div>
                  <div className="bg-purple-500/20 border border-purple-500/30 rounded-lg p-3">
                    <h5 className="text-purple-300 text-sm font-medium">Manager</h5>
                    <p className="text-gray-400 text-xs">Tunnels</p>
                  </div>
                  <div className="bg-purple-500/20 border border-purple-500/30 rounded-lg p-3">
                    <h5 className="text-purple-300 text-sm font-medium">Router</h5>
                    <p className="text-gray-400 text-xs">Subdomains</p>
                  </div>
                </div>
              </motion.div>
            </div>
          </motion.div>
        </div>

        {/* Technical Details */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.4 }}
          viewport={{ once: true }}
          className="mt-16 grid grid-cols-1 md:grid-cols-3 gap-6"
        >
          <motion.div
            whileHover={{ scale: 1.02 }}
            className="bg-dark-800 rounded-xl p-6 border border-gray-700"
          >
            <h3 className="text-lg font-semibold mb-3 text-tunnel-300">Security</h3>
            <ul className="text-gray-300 text-sm space-y-2">
              <li>• TLS 1.3 encryption</li>
              <li>• Token-based authentication</li>
              <li>• Connection tracking</li>
              <li>• Rate limiting support</li>
            </ul>
          </motion.div>

          <motion.div
            whileHover={{ scale: 1.02 }}
            className="bg-dark-800 rounded-xl p-6 border border-gray-700"
          >
            <h3 className="text-lg font-semibold mb-3 text-primary-300">Performance</h3>
            <ul className="text-gray-300 text-sm space-y-2">
              <li>• Go concurrency</li>
              <li>• WebSocket efficiency</li>
              <li>• Minimal memory footprint</li>
              <li>• Connection pooling</li>
            </ul>
          </motion.div>

          <motion.div
            whileHover={{ scale: 1.02 }}
            className="bg-dark-800 rounded-xl p-6 border border-gray-700"
          >
            <h3 className="text-lg font-semibold mb-3 text-purple-300">Deployment</h3>
            <ul className="text-gray-300 text-sm space-y-2">
              <li>• Docker support</li>
              <li>• Health checks</li>
              <li>• Graceful shutdown</li>
              <li>• Configuration management</li>
            </ul>
          </motion.div>
        </motion.div>
      </div>
    </section>
  );
};

export default Architecture; 