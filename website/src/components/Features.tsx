import React from 'react';
import { motion } from 'framer-motion';
import { 
  LockClosedIcon,
  ServerIcon,
  GlobeAltIcon,
  BoltIcon,
  ShieldCheckIcon,
  CogIcon,
  ChartBarIcon,
  CloudIcon
} from '@heroicons/react/24/outline';

const Features: React.FC = () => {
  const features = [
    {
      icon: LockClosedIcon,
      title: "TLS Encryption",
      description: "End-to-end encryption with TLS 1.3 for secure communication between client and server.",
      color: "from-green-500 to-emerald-600"
    },
    {
      icon: ServerIcon,
      title: "Self-Hosted",
      description: "Deploy on your own infrastructure with full control and zero vendor lock-in.",
      color: "from-blue-500 to-cyan-600"
    },
    {
      icon: GlobeAltIcon,
      title: "Subdomain Routing",
      description: "Automatic subdomain mapping for easy access to your local services.",
      color: "from-purple-500 to-pink-600"
    },
    {
      icon: BoltIcon,
      title: "High Performance",
      description: "Built in Go for maximum speed and efficiency with minimal resource usage.",
      color: "from-yellow-500 to-orange-600"
    },
    {
      icon: ShieldCheckIcon,
      title: "Token Authentication",
      description: "Secure token-based authentication with optional public/private key support.",
      color: "from-red-500 to-pink-600"
    },
    {
      icon: CogIcon,
      title: "Easy Configuration",
      description: "Simple YAML/JSON configuration with command-line overrides and sensible defaults.",
      color: "from-indigo-500 to-purple-600"
    },
    {
      icon: ChartBarIcon,
      title: "Connection Tracking",
      description: "Real-time monitoring and connection tracking with detailed logging.",
      color: "from-teal-500 to-green-600"
    },
    {
      icon: CloudIcon,
      title: "Docker Ready",
      description: "Production-ready Docker images with health checks and orchestration support.",
      color: "from-sky-500 to-blue-600"
    }
  ];

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1
      }
    }
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: {
        duration: 0.5
      }
    }
  };

  return (
    <section id="features" className="py-20 bg-dark-800/50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
          className="text-center mb-16"
        >
          <h2 className="text-4xl md:text-5xl font-bold mb-6">
            <span className="text-white">Powerful</span>
            <span className="text-gradient"> Features</span>
          </h2>
          <p className="text-xl text-gray-300 max-w-3xl mx-auto">
            GoTunnel provides enterprise-grade tunneling capabilities with the simplicity 
            and performance you need for modern development workflows.
          </p>
        </motion.div>

        <motion.div
          variants={containerVariants}
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true }}
          className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6"
        >
          {features.map((feature, index) => (
            <motion.div
              key={index}
              variants={itemVariants}
              whileHover={{ scale: 1.05, y: -5 }}
              className="feature-card group"
            >
              <div className={`w-12 h-12 rounded-lg bg-gradient-to-r ${feature.color} flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-200`}>
                <feature.icon className="w-6 h-6 text-white" />
              </div>
              <h3 className="text-lg font-semibold mb-2 text-white group-hover:text-tunnel-300 transition-colors duration-200">
                {feature.title}
              </h3>
              <p className="text-gray-400 text-sm leading-relaxed">
                {feature.description}
              </p>
            </motion.div>
          ))}
        </motion.div>

        {/* Additional highlights */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.4 }}
          viewport={{ once: true }}
          className="mt-16 grid grid-cols-1 lg:grid-cols-3 gap-8"
        >
          <motion.div
            whileHover={{ scale: 1.02 }}
            className="bg-gradient-to-r from-tunnel-500/10 to-primary-500/10 border border-tunnel-500/20 rounded-xl p-6"
          >
            <h3 className="text-xl font-semibold mb-3 text-tunnel-300">HTTP & Raw TCP</h3>
            <p className="text-gray-300">
              Support for both HTTP/HTTPS traffic and raw TCP forwarding for any protocol.
            </p>
          </motion.div>

          <motion.div
            whileHover={{ scale: 1.02 }}
            className="bg-gradient-to-r from-primary-500/10 to-tunnel-500/10 border border-primary-500/20 rounded-xl p-6"
          >
            <h3 className="text-xl font-semibold mb-3 text-primary-300">WebSocket Support</h3>
            <p className="text-gray-300">
              Native WebSocket support for real-time bidirectional communication.
            </p>
          </motion.div>

          <motion.div
            whileHover={{ scale: 1.02 }}
            className="bg-gradient-to-r from-purple-500/10 to-pink-500/10 border border-purple-500/20 rounded-xl p-6"
          >
            <h3 className="text-xl font-semibold mb-3 text-purple-300">Graceful Shutdown</h3>
            <p className="text-gray-300">
              Proper signal handling and graceful shutdown for production deployments.
            </p>
          </motion.div>
        </motion.div>
      </div>
    </section>
  );
};

export default Features; 