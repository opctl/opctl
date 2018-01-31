import nodeFactory from './nodeFactory';
import PeerInfo from 'peer-info';
import * as peerId from 'peer-id';

class LibP2p {
  nodePromise = nodeFactory.construct();

  async dial() {
    const node = await this.nodePromise;
    return new Promise((resolve, reject) => {
      const peer = new PeerInfo(peerId.createFromB58String('QmWX2jb1QTiFViQYayeR6W3sat1BYMcB3iWHFztozLifsk'));
      peer.multiaddrs.add('/ip4/127.0.0.1/tcp/42225/ws');
      node.dial(peer, '/opspec/0.1.5', (err, conn) => {
        if (err) {
          reject(err);
        }
        resolve(conn)
      });
    });
  }
}

export default new LibP2p();
