import peerInfoFactory from './peerInfoFactory';
import Libp2pWebsockets from 'libp2p-websockets';
import libp2pMultiplex from 'libp2p-multiplex';
import libp2pSecio from 'libp2p-secio';
import libp2p from 'libp2p';

class NodeFactory {
  /**
   * constructs (& starts) the node
   * @returns {Promise<void>}
   */
  async construct() {
    const peerInfo = await peerInfoFactory.construct();

    const modules = {
      transport: [new Libp2pWebsockets()],
      connection: {
        muxer: [libp2pMultiplex],
        crypto: [libp2pSecio]
      },
    };

    const node = new libp2p(modules, peerInfo, null, {});
    return new Promise((resolve, reject) => {
      node.start(err => {
        if (err){
          reject(err);
        } else{
          resolve(node);
        }
      })
    });
  }
}

export default new NodeFactory();
