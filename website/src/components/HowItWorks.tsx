import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { 
  ComputerDesktopIcon,
  ServerIcon,
  GlobeAltIcon,
  ArrowRightIcon,
  ArrowLeftIcon
} from '@heroicons/react/24/outline';

const HowItWorks: React.FC = () => {
  const [activeStep, setActiveStep] = useState(0);

  const steps = [
    {
      title: "Client Connects",
      description: "Your GoTunnel client establishes a secure WebSocket connection to the server with authentication.",
      icon: ComputerDesktopIcon,
      color: "from-blue-500 to-cyan-600"
    },
    {
      title: "Server Registers",
      description: "The server registers your client and assigns a subdomain for routing incoming traffic.",
      icon: ServerIcon,
      color: "from-green-500 to-emerald-600"
    },
    {
      title: "Traffic Routing",
      description: "When users visit your subdomain, the server routes traffic through the tunnel to your local service.",
      icon: GlobeAltIcon,
      color: "from-purple-500 to-pink-600"
    }
  ];

  return (
    <section id="how-it-works" className="py-20 bg-dark-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
          className="text-center mb-16"
        >
          <h2 className="text-4xl md:text-5xl font-bold mb-6">
            How <span className="text-gradient">GoTunnel</span> Works
          </h2>
          <p className="text-xl text-gray-300 max-w-3xl mx-auto">
            A simple three-step process to expose your local services securely to the internet.
          </p>
        </motion.div>

        {/* Interactive Steps */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          {/* Steps Navigation */}
          <motion.div
            initial={{ opacity: 0, x: -30 }}
            whileInView={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
            className="space-y-6"
          >
            {steps.map((step, index) => (
              <motion.div
                key={index}
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
                onClick={() => setActiveStep(index)}
                className={`cursor-pointer p-6 rounded-xl border transition-all duration-300 ${
                  activeStep === index
                    ? 'border-tunnel-500 bg-tunnel-500/10'
                    : 'border-gray-700 bg-dark-800/50 hover:border-gray-600'
                }`}
              >
                <div className="flex items-center space-x-4">
                  <div className={`w-12 h-12 rounded-lg bg-gradient-to-r ${step.color} flex items-center justify-center flex-shrink-0`}>
                    <step.icon className="w-6 h-6 text-white" />
                  </div>
                  <div className="flex-1">
                    <h3 className="text-lg font-semibold text-white mb-2">
                      {index + 1}. {step.title}
                    </h3>
                    <p className="text-gray-400 text-sm leading-relaxed">
                      {step.description}
                    </p>
                  </div>
                  {activeStep === index && (
                    <motion.div
                      initial={{ scale: 0 }}
                      animate={{ scale: 1 }}
                      className="w-3 h-3 bg-tunnel-400 rounded-full"
                    />
                  )}
                </div>
              </motion.div>
            ))}
          </motion.div>

          {/* Visual Diagram */}
          <motion.div
            initial={{ opacity: 0, x: 30 }}
            whileInView={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
            className="relative"
          >
            <div className="bg-dark-800 rounded-2xl p-8 border border-gray-700">
              {/* Client */}
              <motion.div
                animate={{
                  x: activeStep === 0 ? [0, 10, 0] : 0,
                  y: activeStep === 0 ? [0, -5, 0] : 0
                }}
                transition={{ duration: 1, repeat: activeStep === 0 ? Infinity : 0 }}
                className="flex items-center space-x-3 mb-8"
              >
                <div className="w-12 h-12 bg-gradient-to-r from-blue-500 to-cyan-600 rounded-lg flex items-center justify-center">
                  <ComputerDesktopIcon className="w-6 h-6 text-white" />
                </div>
                <div>
                  <h4 className="text-white font-semibold">Local Client</h4>
                  <p className="text-gray-400 text-sm">Your development machine</p>
                </div>
              </motion.div>

              {/* Connection Line */}
              <div className="relative mb-8">
                <div className="h-1 bg-gradient-to-r from-blue-500 to-green-500 rounded-full"></div>
                <motion.div
                  animate={{
                    scaleX: activeStep >= 1 ? [0, 1, 0] : 0
                  }}
                  transition={{ duration: 2, repeat: activeStep >= 1 ? Infinity : 0 }}
                  className="absolute top-0 left-0 h-1 bg-tunnel-400 rounded-full origin-left"
                  style={{ width: '100%' }}
                />
              </div>

              {/* Server */}
              <motion.div
                animate={{
                  scale: activeStep === 1 ? [1, 1.1, 1] : 1,
                  boxShadow: activeStep === 1 ? "0 0 20px rgba(14, 165, 233, 0.5)" : "none"
                }}
                transition={{ duration: 1, repeat: activeStep === 1 ? Infinity : 0 }}
                className="flex items-center space-x-3 mb-8"
              >
                <div className="w-12 h-12 bg-gradient-to-r from-green-500 to-emerald-600 rounded-lg flex items-center justify-center">
                  <ServerIcon className="w-6 h-6 text-white" />
                </div>
                <div>
                  <h4 className="text-white font-semibold">GoTunnel Server</h4>
                  <p className="text-gray-400 text-sm">Your hosted infrastructure</p>
                </div>
              </motion.div>

              {/* Connection Line */}
              <div className="relative mb-8">
                <div className="h-1 bg-gradient-to-r from-green-500 to-purple-500 rounded-full"></div>
                <motion.div
                  animate={{
                    scaleX: activeStep >= 2 ? [0, 1, 0] : 0
                  }}
                  transition={{ duration: 2, repeat: activeStep >= 2 ? Infinity : 0 }}
                  className="absolute top-0 left-0 h-1 bg-tunnel-400 rounded-full origin-left"
                  style={{ width: '100%' }}
                />
              </div>

              {/* Internet */}
              <motion.div
                animate={{
                  x: activeStep === 2 ? [0, -5, 0] : 0,
                  y: activeStep === 2 ? [0, 5, 0] : 0
                }}
                transition={{ duration: 1, repeat: activeStep === 2 ? Infinity : 0 }}
                className="flex items-center space-x-3"
              >
                <div className="w-12 h-12 bg-gradient-to-r from-purple-500 to-pink-600 rounded-lg flex items-center justify-center">
                  <GlobeAltIcon className="w-6 h-6 text-white" />
                </div>
                <div>
                  <h4 className="text-white font-semibold">Internet Users</h4>
                  <p className="text-gray-400 text-sm">Access via subdomain</p>
                </div>
              </motion.div>
            </div>

            {/* Data Flow Animation */}
            <motion.div
              animate={{
                opacity: activeStep >= 1 ? [0, 1, 0] : 0
              }}
              transition={{ duration: 2, repeat: activeStep >= 1 ? Infinity : 0 }}
              className="absolute inset-0 flex items-center justify-center pointer-events-none"
            >
              <div className="w-4 h-4 bg-tunnel-400 rounded-full animate-pulse"></div>
            </motion.div>
          </motion.div>
        </div>

        {/* Code Example */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.4 }}
          viewport={{ once: true }}
          className="mt-16"
        >
          <div className="bg-dark-800 rounded-xl p-6 border border-gray-700">
            <h3 className="text-xl font-semibold mb-4 text-white">Quick Setup</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <h4 className="text-tunnel-300 font-semibold mb-2">1. Start the Server</h4>
                <div className="code-block">
                  <pre className="text-green-400">
{`$ ./gotunnel-server \\
  --port 443 \\
  --tls \\
  --cert server.crt \\
  --key server.key \\
  --allowed-tokens "your-secret-token"`}
                  </pre>
                </div>
              </div>
              <div>
                <h4 className="text-tunnel-300 font-semibold mb-2">2. Connect Client</h4>
                <div className="code-block">
                  <pre className="text-green-400">
{`$ ./gotunnel-client \\
  --server your-server.com \\
  --subdomain myapp \\
  --local-port 3000 \\
  --token "your-secret-token"`}
                  </pre>
                </div>
              </div>
            </div>
            <div className="mt-6 p-4 bg-tunnel-500/10 border border-tunnel-500/20 rounded-lg">
              <p className="text-tunnel-300 text-sm">
                <strong>Result:</strong> Your local service on port 3000 is now accessible at{' '}
                <code className="bg-dark-700 px-2 py-1 rounded">https://myapp.your-server.com</code>
              </p>
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  );
};

export default HowItWorks; 