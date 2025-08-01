import React from 'react';
import { motion } from 'framer-motion';
import { TypeAnimation } from 'react-type-animation';
import { 
  ServerIcon, 
  ShieldCheckIcon, 
  BoltIcon,
  ArrowRightIcon 
} from '@heroicons/react/24/outline';

const Hero: React.FC = () => {
  return (
    <section className="relative min-h-screen flex items-center justify-center overflow-hidden">
      {/* Animated background */}
      <div className="absolute inset-0 tunnel-gradient opacity-20"></div>
      <div className="absolute inset-0 bg-dark-900/50"></div>
      
      {/* Floating particles */}
      <div className="absolute inset-0 overflow-hidden">
        {[...Array(20)].map((_, i) => (
          <motion.div
            key={i}
            className="absolute w-2 h-2 bg-tunnel-400 rounded-full opacity-30"
            animate={{
              x: [0, Math.random() * window.innerWidth],
              y: [0, Math.random() * window.innerHeight],
            }}
            transition={{
              duration: Math.random() * 10 + 10,
              repeat: Infinity,
              ease: "linear"
            }}
            style={{
              left: Math.random() * 100 + '%',
              top: Math.random() * 100 + '%',
            }}
          />
        ))}
      </div>

      <div className="relative z-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          className="mb-8"
        >
          <motion.div
            whileHover={{ scale: 1.05 }}
            className="inline-flex items-center px-4 py-2 bg-dark-800/50 backdrop-blur-sm rounded-full border border-tunnel-500/30 mb-8"
          >
            <ServerIcon className="w-5 h-5 text-tunnel-400 mr-2" />
            <span className="text-tunnel-300 text-sm font-medium">
              Self-Hosted Secure Tunneling
            </span>
          </motion.div>
        </motion.div>

        <motion.h1
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.2 }}
          className="text-5xl md:text-7xl font-bold mb-6"
        >
          <span className="text-white">GoTunnel</span>
          <br />
          <span className="text-gradient">
            <TypeAnimation
              sequence={[
                'Secure Tunneling',
                2000,
                'ngrok Alternative',
                2000,
                'Cloudflare Tunnel',
                2000,
                'Self-Hosted',
                2000,
              ]}
              wrapper="span"
              speed={50}
              repeat={Infinity}
            />
          </span>
        </motion.h1>

        <motion.p
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.4 }}
          className="text-xl md:text-2xl text-gray-300 mb-8 max-w-4xl mx-auto leading-relaxed"
        >
          Build your own secure tunneling infrastructure with Go. 
          Expose local services to the internet with enterprise-grade security, 
          full control, and zero vendor lock-in.
        </motion.p>

        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.6 }}
          className="flex flex-col sm:flex-row gap-4 justify-center items-center mb-12"
        >
          <motion.button
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            className="bg-gradient-to-r from-tunnel-500 to-primary-600 text-white px-8 py-4 rounded-lg font-semibold text-lg hover:shadow-xl transition-all duration-200 flex items-center"
          >
            Get Started
            <ArrowRightIcon className="w-5 h-5 ml-2" />
          </motion.button>
          
          <motion.button
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            className="border border-tunnel-500 text-tunnel-300 px-8 py-4 rounded-lg font-semibold text-lg hover:bg-tunnel-500/10 transition-all duration-200"
          >
            View Documentation
          </motion.button>
        </motion.div>

        {/* Feature highlights */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.8 }}
          className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-4xl mx-auto"
        >
          <motion.div
            whileHover={{ scale: 1.05 }}
            className="feature-card text-center"
          >
            <ShieldCheckIcon className="w-12 h-12 text-tunnel-400 mx-auto mb-4" />
            <h3 className="text-lg font-semibold mb-2">Enterprise Security</h3>
            <p className="text-gray-400">TLS encryption, token-based auth, and full control over your infrastructure</p>
          </motion.div>

          <motion.div
            whileHover={{ scale: 1.05 }}
            className="feature-card text-center"
          >
            <BoltIcon className="w-12 h-12 text-tunnel-400 mx-auto mb-4" />
            <h3 className="text-lg font-semibold mb-2">High Performance</h3>
            <p className="text-gray-400">Built in Go for maximum speed and efficiency with minimal resource usage</p>
          </motion.div>

          <motion.div
            whileHover={{ scale: 1.05 }}
            className="feature-card text-center"
          >
            <ServerIcon className="w-12 h-12 text-tunnel-400 mx-auto mb-4" />
            <h3 className="text-lg font-semibold mb-2">Self-Hosted</h3>
            <p className="text-gray-400">Deploy on your own infrastructure with zero vendor lock-in or limits</p>
          </motion.div>
        </motion.div>
      </div>

      {/* Scroll indicator */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 1.5 }}
        className="absolute bottom-8 left-1/2 transform -translate-x-1/2"
      >
        <motion.div
          animate={{ y: [0, 10, 0] }}
          transition={{ duration: 2, repeat: Infinity }}
          className="w-6 h-10 border-2 border-tunnel-400 rounded-full flex justify-center"
        >
          <motion.div
            animate={{ y: [0, 12, 0] }}
            transition={{ duration: 2, repeat: Infinity }}
            className="w-1 h-3 bg-tunnel-400 rounded-full mt-2"
          />
        </motion.div>
      </motion.div>
    </section>
  );
};

export default Hero; 