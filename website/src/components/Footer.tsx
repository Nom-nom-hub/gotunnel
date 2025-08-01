import React from 'react';
import { motion } from 'framer-motion';
import { 
  GithubIcon,
  TwitterIcon,
  LinkedinIcon,
  MailIcon,
  HeartIcon
} from '@heroicons/react/24/outline';

const Footer: React.FC = () => {
  const currentYear = new Date().getFullYear();

  const links = {
    product: [
      { name: 'Features', href: '#features' },
      { name: 'How It Works', href: '#how-it-works' },
      { name: 'Quick Start', href: '#quick-start' },
      { name: 'Examples', href: '#examples' }
    ],
    resources: [
      { name: 'Documentation', href: '#' },
      { name: 'API Reference', href: '#' },
      { name: 'Tutorials', href: '#' },
      { name: 'Blog', href: '#' }
    ],
    community: [
      { name: 'GitHub', href: 'https://github.com/ogrok/gotunnel' },
      { name: 'Discussions', href: '#' },
      { name: 'Contributing', href: '#' },
      { name: 'Issues', href: '#' }
    ],
    company: [
      { name: 'About', href: '#' },
      { name: 'Privacy', href: '#' },
      { name: 'Terms', href: '#' },
      { name: 'Contact', href: '#' }
    ]
  };

  const socialLinks = [
    { name: 'GitHub', icon: GithubIcon, href: 'https://github.com/ogrok/gotunnel' },
    { name: 'Twitter', icon: TwitterIcon, href: '#' },
    { name: 'LinkedIn', icon: LinkedinIcon, href: '#' },
    { name: 'Email', icon: MailIcon, href: 'mailto:contact@gotunnel.dev' }
  ];

  return (
    <footer className="bg-dark-900 border-t border-gray-800">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-8">
          {/* Brand */}
          <div className="lg:col-span-2">
            <motion.div
              whileHover={{ scale: 1.05 }}
              className="flex items-center space-x-2 mb-4"
            >
              <div className="w-8 h-8 bg-gradient-to-r from-tunnel-400 to-primary-500 rounded-lg flex items-center justify-center">
                <span className="text-white font-bold text-sm">GT</span>
              </div>
              <span className="text-xl font-bold text-gradient">GoTunnel</span>
            </motion.div>
            
            <p className="text-gray-400 mb-6 max-w-md">
              Self-hosted secure tunneling system written in Go. 
              Expose local services to the internet with enterprise-grade security, 
              full control, and zero vendor lock-in.
            </p>

            {/* Social Links */}
            <div className="flex space-x-4">
              {socialLinks.map((social) => (
                <motion.a
                  key={social.name}
                  href={social.href}
                  target="_blank"
                  rel="noopener noreferrer"
                  whileHover={{ scale: 1.1, y: -2 }}
                  whileTap={{ scale: 0.95 }}
                  className="w-10 h-10 bg-dark-800 border border-gray-700 rounded-lg flex items-center justify-center text-gray-400 hover:text-tunnel-400 hover:border-tunnel-500/50 transition-all duration-200"
                >
                  <social.icon className="w-5 h-5" />
                </motion.a>
              ))}
            </div>
          </div>

          {/* Product Links */}
          <div>
            <h3 className="text-white font-semibold mb-4">Product</h3>
            <ul className="space-y-2">
              {links.product.map((link) => (
                <li key={link.name}>
                  <a
                    href={link.href}
                    className="text-gray-400 hover:text-tunnel-300 transition-colors duration-200 text-sm"
                  >
                    {link.name}
                  </a>
                </li>
              ))}
            </ul>
          </div>

          {/* Resources Links */}
          <div>
            <h3 className="text-white font-semibold mb-4">Resources</h3>
            <ul className="space-y-2">
              {links.resources.map((link) => (
                <li key={link.name}>
                  <a
                    href={link.href}
                    className="text-gray-400 hover:text-tunnel-300 transition-colors duration-200 text-sm"
                  >
                    {link.name}
                  </a>
                </li>
              ))}
            </ul>
          </div>

          {/* Community Links */}
          <div>
            <h3 className="text-white font-semibold mb-4">Community</h3>
            <ul className="space-y-2">
              {links.community.map((link) => (
                <li key={link.name}>
                  <a
                    href={link.href}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-gray-400 hover:text-tunnel-300 transition-colors duration-200 text-sm"
                  >
                    {link.name}
                  </a>
                </li>
              ))}
            </ul>
          </div>
        </div>

        {/* Bottom Section */}
        <div className="mt-12 pt-8 border-t border-gray-800">
          <div className="flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
            <div className="flex items-center space-x-2 text-gray-400 text-sm">
              <span>© {currentYear} GoTunnel. Made with</span>
              <motion.div
                animate={{ scale: [1, 1.2, 1] }}
                transition={{ duration: 1, repeat: Infinity }}
              >
                <HeartIcon className="w-4 h-4 text-red-500" />
              </motion.div>
              <span>by the GoTunnel team.</span>
            </div>

            <div className="flex items-center space-x-6 text-sm">
              <a
                href="#"
                className="text-gray-400 hover:text-tunnel-300 transition-colors duration-200"
              >
                Privacy Policy
              </a>
              <a
                href="#"
                className="text-gray-400 hover:text-tunnel-300 transition-colors duration-200"
              >
                Terms of Service
              </a>
              <a
                href="#"
                className="text-gray-400 hover:text-tunnel-300 transition-colors duration-200"
              >
                Cookie Policy
              </a>
            </div>
          </div>

          {/* Newsletter Signup */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
            className="mt-8 p-6 bg-gradient-to-r from-tunnel-500/10 to-primary-500/10 border border-tunnel-500/20 rounded-xl"
          >
            <div className="text-center">
              <h3 className="text-lg font-semibold text-white mb-2">
                Stay Updated
              </h3>
              <p className="text-gray-300 mb-4">
                Get the latest updates, tutorials, and community news.
              </p>
              
              <div className="flex flex-col sm:flex-row gap-3 max-w-md mx-auto">
                <input
                  type="email"
                  placeholder="Enter your email"
                  className="flex-1 px-4 py-2 bg-dark-800 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:border-tunnel-500 transition-colors duration-200"
                />
                <motion.button
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                  className="px-6 py-2 bg-gradient-to-r from-tunnel-500 to-primary-600 text-white rounded-lg font-medium hover:shadow-lg transition-all duration-200"
                >
                  Subscribe
                </motion.button>
              </div>
            </div>
          </motion.div>
        </div>
      </div>
    </footer>
  );
};

export default Footer; 