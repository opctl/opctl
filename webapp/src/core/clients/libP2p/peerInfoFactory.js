import PeerInfo from 'peer-info';

class PeerInfoFactory {
  /**
   * constructs peer info for the node
   * @returns {Promise<any>}
   */
  async construct() {
    return new Promise((resolve, reject) => {
      PeerInfo.create((err, peerInfo) => {
        if (err) {
          reject(err);
        } else {
          resolve(peerInfo);
        }
      });
    });
  }
}

export default new PeerInfoFactory();
