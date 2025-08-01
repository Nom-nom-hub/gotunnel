import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { 
  PlayIcon,
  CommandLineIcon,
  ServerIcon,
  ComputerDesktopIcon,
  CheckCircleIcon
} from '@heroicons/react/24/outline';

const QuickStart: React.FC = () => {
  const [activeStep, setActiveStep] = useState(0);

  const steps = [
    {
      title: "Install GoTunnel",
      description: "Download and build the GoTunnel binaries",
      icon: CommandLineIcon,
      commands: [
        "git clone https://github.com/ogrok/gotunnel.git",
        "cd gotunnel",
        "make build"
      ],
      output: "✓ Built gotunnel-server and gotunnel-client"
    },
    {
      title: "Generate Certificates",
      description: "Create self-signed TLS certificates for local development",
      icon: ServerIcon,
      commands: [
        "make certs",
        "ls -la certs/"
      ],
      output: "✓ Generated server.crt and server.key"
    },
    {
      title: "Start the Server",
      description: "Launch the GoTunnel server with TLS enabled",
      icon: ServerIcon,
      commands: [
        "./gotunnel-server \\",
        "  --port 443 \\",
        "  --tls \\",
        "  --cert certs/server.crt \\",
        "  --key certs/server.key \\",
        "  --allowed-tokens 'your-secret-token'"
      ],
      output: "✓ Server started on :443 with TLS"
    },
    {
      title: "Connect Client",
      description: "Connect your local service to the tunnel",
      icon: ComputerDesktopIcon,
      commands: [
        "./gotunnel-client \\",
        "  --server localhost \\",
        "  --subdomain myapp \\",
        "  --local-port 3000 \\",
        "  --token 'your-secret-token'"
      ],
      output: "✓ Connected! Access at https://myapp.localhost"
    }
  ];

  return (
    <section id="quick-start" className="py-20 bg-dark-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
          className="text-center mb-16"
        >
          <h2 className="text-4xl md:text-5xl font-bold mb-6">
            <span className="text-white">Quick</span>
            <span className="text-gradient"> Start</span>
          </h2>
          <p className="text-xl text-gray-300 max-w-3xl mx-auto">
            Get GoTunnel running in minutes with our step-by-step guide.
          </p>
        </motion.div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-start">
          {/* Step Navigation */}
          <motion.div
            initial={{ opacity: 0, x: -30 }}
            whileInView={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
            className="space-y-4"
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
                  <div className={`w-12 h-12 rounded-lg bg-gradient-to-r ${
                    activeStep === index ? 'from-tunnel-500 to-primary-600' : 'from-gray-600 to-gray-700'
                  } flex items-center justify-center flex-shrink-0`}>
                    <step.icon className="w-6 h-6 text-white" />
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center space-x-2 mb-2">
                      <span className="text-sm font-medium text-gray-400">Step {index + 1}</span>
                      {activeStep === index && (
                        <motion.div
                          initial={{ scale: 0 }}
                          animate={{ scale: 1 }}
                          className="w-2 h-2 bg-tunnel-400 rounded-full"
                        />
                      )}
                    </div>
                    <h3 className="text-lg font-semibold text-white mb-2">
                      {step.title}
                    </h3>
                    <p className="text-gray-400 text-sm leading-relaxed">
                      {step.description}
                    </p>
                  </div>
                </div>
              </motion.div>
            ))}
          </motion.div>

          {/* Terminal Output */}
          <motion.div
            initial={{ opacity: 0, x: 30 }}
            whileInView={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
            className="relative"
          >
            <div className="terminal-window">
              <div className="terminal-header">
                <div className="terminal-dot bg-red-500"></div>
                <div className="terminal-dot bg-yellow-500"></div>
                <div className="terminal-dot bg-green-500"></div>
                <span className="text-gray-300 text-sm ml-4">GoTunnel Terminal</span>
              </div>
              <div className="terminal-content">
                <motion.div
                  key={activeStep}
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.3 }}
                >
                  <div className="mb-4">
                    <span className="text-green-400">$</span>
                    <span className="text-white ml-2">GoTunnel Setup</span>
                  </div>
                  
                  {steps[activeStep].commands.map((command, index) => (
                    <div key={index} className="mb-2">
                      <div className="flex items-center">
                        <span className="text-green-400">$</span>
                        <span className="text-white ml-2 font-mono text-sm">{command}</span>
                      </div>
                    </div>
                  ))}
                  
                  <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ delay: 0.5 }}
                    className="mt-4 p-3 bg-green-500/10 border border-green-500/20 rounded-lg"
                  >
                    <div className="flex items-center space-x-2">
                      <CheckCircleIcon className="w-4 h-4 text-green-400" />
                      <span className="text-green-400 text-sm">{steps[activeStep].output}</span>
                    </div>
                  </motion.div>
                </motion.div>
              </div>
            </div>

            {/* Progress Indicator */}
            <div className="mt-6 flex justify-center">
              <div className="flex space-x-2">
                {steps.map((_, index) => (
                  <motion.div
                    key={index}
                    className={`w-3 h-3 rounded-full transition-all duration-300 ${
                      index <= activeStep ? 'bg-tunnel-400' : 'bg-gray-600'
                    }`}
                    whileHover={{ scale: 1.2 }}
                  />
                ))}
              </div>
            </div>
          </motion.div>
        </div>

        {/* Configuration Example */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.4 }}
          viewport={{ once: true }}
          className="mt-16"
        >
          <div className="bg-dark-800 rounded-xl p-6 border border-gray-700">
            <h3 className="text-xl font-semibold mb-4 text-white">Configuration File</h3>
            <p className="text-gray-300 mb-4">Create a <code className="bg-dark-700 px-2 py-1 rounded">config.yaml</code> file for easier client configuration:</p>
            <div className="code-block">
              <pre className="text-gray-300">
{`# config.yaml
server_addr: "your-server.com"
subdomain: "myapp"
local_port: 3000
local_host: "localhost"
auth_token: "your-secret-token"
use_tls: true
skip_verify: false`}
              </pre>
            </div>
            <div className="mt-4 p-4 bg-tunnel-500/10 border border-tunnel-500/20 rounded-lg">
              <p className="text-tunnel-300 text-sm">
                <strong>Usage:</strong> <code className="bg-dark-700 px-2 py-1 rounded">./gotunnel-client --config config.yaml</code>
              </p>
            </div>
          </div>
        </motion.div>

        {/* Docker Quick Start */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.6 }}
          viewport={{ once: true }}
          className="mt-8"
        >
          <div className="bg-dark-800 rounded-xl p-6 border border-gray-700">
            <h3 className="text-xl font-semibold mb-4 text-white">Docker Quick Start</h3>
            <p className="text-gray-300 mb-4">Deploy with Docker for production-ready setup:</p>
            <div className="code-block">
              <pre className="text-gray-300">
{`# Start server with Docker
docker run -d \\
  --name gotunnel-server \\
  -p 443:443 \\
  -v $(pwd)/certs:/certs \\
  -e ALLOWED_TOKENS="your-secret-token" \\
  gotunnel-server

# Or use docker-compose
docker-compose up -d`}
              </pre>
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  );
};

export default QuickStart; 