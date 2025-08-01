import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { 
  CheckIcon,
  XMarkIcon,
  MinusIcon,
  StarIcon
} from '@heroicons/react/24/outline';

const Comparison: React.FC = () => {
  const [activeFeature, setActiveFeature] = useState(0);

  const features = [
    "Self-Hosted",
    "TLS Encryption",
    "WebSocket Support",
    "HTTP/TCP Forwarding",
    "Subdomain Routing",
    "Token Authentication",
    "Docker Support",
    "Rate Limiting",
    "Connection Tracking",
    "Graceful Shutdown",
    "Custom Domains",
    "Zero Vendor Lock-in"
  ];

  const solutions = [
    {
      name: "GoTunnel",
      description: "Self-hosted secure tunneling",
      logo: "GT",
      color: "from-tunnel-500 to-primary-600",
      features: [true, true, true, true, true, true, true, true, true, true, true, true],
      pros: [
        "Complete control over infrastructure",
        "No usage limits or costs",
        "Enterprise-grade security",
        "Full source code access"
      ],
      cons: [
        "Requires server setup",
        "Manual certificate management"
      ],
      rating: 5
    },
    {
      name: "ngrok",
      description: "Popular tunneling service",
      logo: "ng",
      color: "from-green-500 to-emerald-600",
      features: [false, true, true, true, true, true, false, true, true, true, true, false],
      pros: [
        "Easy to use",
        "Free tier available",
        "Good documentation",
        "Reliable service"
      ],
      cons: [
        "Usage limits on free tier",
        "Vendor lock-in",
        "Limited customization",
        "Privacy concerns"
      ],
      rating: 4
    },
    {
      name: "Cloudflare Tunnel",
      description: "Cloudflare's tunneling solution",
      logo: "CF",
      color: "from-orange-500 to-red-600",
      features: [false, true, false, true, true, true, true, true, true, true, true, false],
      pros: [
        "Integrated with Cloudflare",
        "Good security features",
        "DDoS protection",
        "Global CDN"
      ],
      cons: [
        "Requires Cloudflare account",
        "Limited to HTTP/HTTPS",
        "Vendor lock-in",
        "Complex setup"
      ],
      rating: 4
    }
  ];

  const getFeatureIcon = (supported: boolean) => {
    if (supported) {
      return <CheckIcon className="w-5 h-5 text-green-400" />;
    }
    return <XMarkIcon className="w-5 h-5 text-red-400" />;
  };

  return (
    <section className="py-20 bg-dark-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
          className="text-center mb-16"
        >
          <h2 className="text-4xl md:text-5xl font-bold mb-6">
            <span className="text-white">Compare</span>
            <span className="text-gradient"> Solutions</span>
          </h2>
          <p className="text-xl text-gray-300 max-w-3xl mx-auto">
            See how GoTunnel stacks up against popular tunneling solutions.
          </p>
        </motion.div>

        {/* Feature Comparison Table */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.2 }}
          viewport={{ once: true }}
          className="overflow-x-auto"
        >
          <div className="bg-dark-800 rounded-xl border border-gray-700 overflow-hidden">
            <div className="grid grid-cols-4 gap-0">
              {/* Header */}
              <div className="bg-dark-700 p-6 border-r border-gray-700">
                <h3 className="text-lg font-semibold text-white mb-2">Features</h3>
                <p className="text-gray-400 text-sm">Capabilities comparison</p>
              </div>
              
              {solutions.map((solution, index) => (
                <div key={index} className="bg-dark-700 p-6 border-r border-gray-700 last:border-r-0">
                  <div className="flex items-center space-x-3 mb-2">
                    <div className={`w-8 h-8 rounded-lg bg-gradient-to-r ${solution.color} flex items-center justify-center`}>
                      <span className="text-white font-bold text-sm">{solution.logo}</span>
                    </div>
                    <div>
                      <h4 className="text-white font-semibold">{solution.name}</h4>
                      <p className="text-gray-400 text-xs">{solution.description}</p>
                    </div>
                  </div>
                  <div className="flex items-center space-x-1">
                    {[...Array(5)].map((_, i) => (
                      <StarIcon
                        key={i}
                        className={`w-4 h-4 ${
                          i < solution.rating ? 'text-yellow-400 fill-current' : 'text-gray-600'
                        }`}
                      />
                    ))}
                  </div>
                </div>
              ))}
            </div>

            {/* Feature Rows */}
            {features.map((feature, featureIndex) => (
              <motion.div
                key={featureIndex}
                initial={{ opacity: 0, x: -20 }}
                whileInView={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.5, delay: featureIndex * 0.1 }}
                viewport={{ once: true }}
                className={`grid grid-cols-4 gap-0 border-t border-gray-700 ${
                  activeFeature === featureIndex ? 'bg-tunnel-500/5' : ''
                }`}
                onMouseEnter={() => setActiveFeature(featureIndex)}
              >
                <div className="p-4 border-r border-gray-700">
                  <span className="text-gray-300 text-sm font-medium">{feature}</span>
                </div>
                
                {solutions.map((solution, solutionIndex) => (
                  <div key={solutionIndex} className="p-4 border-r border-gray-700 last:border-r-0 flex items-center justify-center">
                    {getFeatureIcon(solution.features[featureIndex])}
                  </div>
                ))}
              </motion.div>
            ))}
          </div>
        </motion.div>

        {/* Detailed Comparison Cards */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.4 }}
          viewport={{ once: true }}
          className="mt-16 grid grid-cols-1 md:grid-cols-3 gap-6"
        >
          {solutions.map((solution, index) => (
            <motion.div
              key={index}
              whileHover={{ scale: 1.02 }}
              className={`bg-dark-800 rounded-xl p-6 border ${
                index === 0 ? 'border-tunnel-500/50 bg-tunnel-500/5' : 'border-gray-700'
              }`}
            >
              <div className="flex items-center space-x-3 mb-4">
                <div className={`w-10 h-10 rounded-lg bg-gradient-to-r ${solution.color} flex items-center justify-center`}>
                  <span className="text-white font-bold">{solution.logo}</span>
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-white">{solution.name}</h3>
                  <p className="text-gray-400 text-sm">{solution.description}</p>
                </div>
              </div>

              {/* Rating */}
              <div className="flex items-center space-x-2 mb-4">
                <div className="flex space-x-1">
                  {[...Array(5)].map((_, i) => (
                    <StarIcon
                      key={i}
                      className={`w-4 h-4 ${
                        i < solution.rating ? 'text-yellow-400 fill-current' : 'text-gray-600'
                      }`}
                    />
                  ))}
                </div>
                <span className="text-gray-400 text-sm">{solution.rating}/5</span>
              </div>

              {/* Pros */}
              <div className="mb-4">
                <h4 className="text-green-400 font-semibold mb-2">Pros</h4>
                <ul className="space-y-1">
                  {solution.pros.map((pro, proIndex) => (
                    <li key={proIndex} className="flex items-start space-x-2 text-gray-300 text-sm">
                      <CheckIcon className="w-4 h-4 text-green-400 mt-0.5 flex-shrink-0" />
                      <span>{pro}</span>
                    </li>
                  ))}
                </ul>
              </div>

              {/* Cons */}
              <div>
                <h4 className="text-red-400 font-semibold mb-2">Cons</h4>
                <ul className="space-y-1">
                  {solution.cons.map((con, conIndex) => (
                    <li key={conIndex} className="flex items-start space-x-2 text-gray-300 text-sm">
                      <XMarkIcon className="w-4 h-4 text-red-400 mt-0.5 flex-shrink-0" />
                      <span>{con}</span>
                    </li>
                  ))}
                </ul>
              </div>
            </motion.div>
          ))}
        </motion.div>

        {/* Why Choose GoTunnel */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.6 }}
          viewport={{ once: true }}
          className="mt-16 bg-gradient-to-r from-tunnel-500/10 to-primary-500/10 border border-tunnel-500/20 rounded-xl p-8"
        >
          <div className="text-center">
            <h3 className="text-2xl font-bold text-white mb-4">Why Choose GoTunnel?</h3>
            <p className="text-gray-300 mb-6 max-w-3xl mx-auto">
              GoTunnel gives you complete control over your tunneling infrastructure while providing 
              enterprise-grade security and performance. No limits, no vendor lock-in, just pure freedom.
            </p>
            
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <motion.div
                whileHover={{ scale: 1.05 }}
                className="text-center"
              >
                <div className="w-12 h-12 bg-gradient-to-r from-green-500 to-emerald-600 rounded-lg flex items-center justify-center mx-auto mb-3">
                  <CheckIcon className="w-6 h-6 text-white" />
                </div>
                <h4 className="text-white font-semibold mb-2">Complete Control</h4>
                <p className="text-gray-400 text-sm">Deploy on your infrastructure with full customization</p>
              </motion.div>

              <motion.div
                whileHover={{ scale: 1.05 }}
                className="text-center"
              >
                <div className="w-12 h-12 bg-gradient-to-r from-blue-500 to-cyan-600 rounded-lg flex items-center justify-center mx-auto mb-3">
                  <StarIcon className="w-6 h-6 text-white" />
                </div>
                <h4 className="text-white font-semibold mb-2">No Limits</h4>
                <p className="text-gray-400 text-sm">Unlimited usage without vendor restrictions</p>
              </motion.div>

              <motion.div
                whileHover={{ scale: 1.05 }}
                className="text-center"
              >
                <div className="w-12 h-12 bg-gradient-to-r from-purple-500 to-pink-600 rounded-lg flex items-center justify-center mx-auto mb-3">
                  <CheckIcon className="w-6 h-6 text-white" />
                </div>
                <h4 className="text-white font-semibold mb-2">Open Source</h4>
                <p className="text-gray-400 text-sm">Full source code access and community-driven</p>
              </motion.div>
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  );
};

export default Comparison; 