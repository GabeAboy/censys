import React, { useState } from 'react';

const CreateAssetModal = ({ isOpen, onClose, onAssetCreated }) => {
  const [formData, setFormData] = useState({
    hostname: '',
    ipAddress: '',
    ports: '',
    tags: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [validationErrors, setValidationErrors] = useState({});

  if (!isOpen) return null;

  const validateForm = () => {
    const errors = {};

    if (!formData.hostname.trim()) {
      errors.hostname = 'Hostname is required';
    }

    if (!formData.ipAddress.trim()) {
      errors.ipAddress = 'IP Address is required';
    } else {
      const ipPattern = /^(\d{1,3}\.){3}\d{1,3}$/;
      if (!ipPattern.test(formData.ipAddress.trim())) {
        errors.ipAddress = 'Invalid IP address format';
      }
    }

    if (!formData.ports.trim()) {
      errors.ports = 'At least one port is required';
    } else {
      const portNumbers = formData.ports.split(',').map(p => parseInt(p.trim())).filter(p => !isNaN(p));
      if (portNumbers.length === 0) {
        errors.ports = 'At least one valid port number is required';
      } else if (portNumbers.some(p => p < 1 || p > 65535)) {
        errors.ports = 'Port numbers must be between 1 and 65535';
      }
    }

    return errors;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);
    setValidationErrors({});

    const errors = validateForm();
    if (Object.keys(errors).length > 0) {
      setValidationErrors(errors);
      return;
    }

    setLoading(true);

    try {
      // Parse ports and tags from comma-separated strings
      const portNumbers = formData.ports
        ? formData.ports.split(',').map(p => parseInt(p.trim())).filter(p => !isNaN(p))
        : [];
      
      const tags = formData.tags
        ? formData.tags.split(',').map(t => t.trim()).filter(t => t.length > 0)
        : [];

      const payload = {
        hostname: formData.hostname,
        ip_address: formData.ipAddress,
        port_numbers: portNumbers,
        tags: tags
      };

      const response = await fetch('http://localhost:8080/api/v1/assets', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to create asset');
      }

      const newAsset = await response.json();

      // Reset form
      setFormData({
        hostname: '',
        ipAddress: '',
        ports: '',
        tags: ''
      });
      setValidationErrors({});

      // Notify parent component
      if (onAssetCreated) {
        onAssetCreated(newAsset);
      }

      onClose();
    } catch (err) {
      console.error('Error creating asset:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">

      <div
        className="fixed inset-0 bg-black bg-opacity-50 transition-opacity"
        onClick={onClose}
      />

      <div className="flex min-h-full items-center justify-center p-4">
        <div className="relative bg-white rounded-lg shadow-xl max-w-md w-full">
          <div className="flex items-center justify-between p-6 border-b border-gray-200">
            <h2 className="text-xl font-semibold text-gray-900">Create New Asset</h2>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-600 transition-colors"
            >
              <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit} className="p-6">
            {error && (
              <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg">
                <p className="text-sm text-red-800">{error}</p>
              </div>
            )}

            <div className="space-y-4">
              <div>
                <label htmlFor="hostname" className="block text-sm font-medium text-gray-700 mb-1">
                  Hostname *
                </label>
                <input
                  type="text"
                  id="hostname"
                  name="hostname"
                  value={formData.hostname}
                  onChange={handleChange}
                  placeholder="web-server-01.example.com"
                  className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none ${
                    validationErrors.hostname ? 'border-red-300 bg-red-50' : 'border-gray-300'
                  }`}
                />
                {validationErrors.hostname && (
                  <p className="mt-1 text-xs text-red-600">{validationErrors.hostname}</p>
                )}
              </div>

              <div>
                <label htmlFor="ipAddress" className="block text-sm font-medium text-gray-700 mb-1">
                  IP Address *
                </label>
                <input
                  type="text"
                  id="ipAddress"
                  name="ipAddress"
                  value={formData.ipAddress}
                  onChange={handleChange}
                  placeholder="192.168.1.100"
                  className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none font-mono ${
                    validationErrors.ipAddress ? 'border-red-300 bg-red-50' : 'border-gray-300'
                  }`}
                />
                {validationErrors.ipAddress && (
                  <p className="mt-1 text-xs text-red-600">{validationErrors.ipAddress}</p>
                )}
              </div>

              <div>
                <label htmlFor="ports" className="block text-sm font-medium text-gray-700 mb-1">
                  Open Ports *
                </label>
                <input
                  type="text"
                  id="ports"
                  name="ports"
                  value={formData.ports}
                  onChange={handleChange}
                  placeholder="22, 80, 443"
                  className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none ${
                    validationErrors.ports ? 'border-red-300 bg-red-50' : 'border-gray-300'
                  }`}
                />
                {validationErrors.ports ? (
                  <p className="mt-1 text-xs text-red-600">{validationErrors.ports}</p>
                ) : (
                  <p className="mt-1 text-xs text-gray-500">Comma-separated port numbers</p>
                )}
              </div>

              <div>
                <label htmlFor="tags" className="block text-sm font-medium text-gray-700 mb-1">
                  Tags
                </label>
                <input
                  type="text"
                  id="tags"
                  name="tags"
                  value={formData.tags}
                  onChange={handleChange}
                  placeholder="production, web, critical"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                />
                <p className="mt-1 text-xs text-gray-500">Comma-separated tags</p>
              </div>

              <div className="bg-blue-50 border border-blue-200 rounded-lg p-3">
                <p className="text-xs text-blue-800">
                  <strong>Risk Level</strong> will be automatically calculated based on open ports:
                </p>
                <ul className="mt-2 text-xs text-blue-700 space-y-1 ml-4 list-disc">
                  <li>High: SSH (22), RDP (3389), FTP (21)</li>
                  <li>Medium: HTTPS (443) with expired cert</li>
                  <li>Low: Other ports</li>
                </ul>
              </div>
            </div>

            <div className="flex gap-3 mt-6">
              <button
                type="button"
                onClick={onClose}
                className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors font-medium"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={loading}
                className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium disabled:bg-blue-400 disabled:cursor-not-allowed"
              >
                {loading ? 'Creating...' : 'Create Asset'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default CreateAssetModal;
