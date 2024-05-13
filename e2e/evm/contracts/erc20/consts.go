// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: LGPL-3.0-only

package erc20

const ERC20PresetMinterPauserABI = `[{"inputs":[{"internalType":"string","name":"name","type":"string"},{"internalType":"string","name":"symbol","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Paused","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"role","type":"bytes32"},{"indexed":true,"internalType":"bytes32","name":"previousAdminRole","type":"bytes32"},{"indexed":true,"internalType":"bytes32","name":"newAdminRole","type":"bytes32"}],"name":"RoleAdminChanged","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"role","type":"bytes32"},{"indexed":true,"internalType":"address","name":"account","type":"address"},{"indexed":true,"internalType":"address","name":"sender","type":"address"}],"name":"RoleGranted","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"role","type":"bytes32"},{"indexed":true,"internalType":"address","name":"account","type":"address"},{"indexed":true,"internalType":"address","name":"sender","type":"address"}],"name":"RoleRevoked","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Unpaused","type":"event"},{"inputs":[],"name":"DEFAULT_ADMIN_ROLE","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MINTER_ROLE","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"PAUSER_ROLE","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"burn","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"burnFrom","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"}],"name":"getRoleAdmin","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"},{"internalType":"uint256","name":"index","type":"uint256"}],"name":"getRoleMember","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"}],"name":"getRoleMemberCount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"},{"internalType":"address","name":"account","type":"address"}],"name":"grantRole","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"},{"internalType":"address","name":"account","type":"address"}],"name":"hasRole","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"paused","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"},{"internalType":"address","name":"account","type":"address"}],"name":"renounceRole","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"},{"internalType":"address","name":"account","type":"address"}],"name":"revokeRole","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes4","name":"interfaceId","type":"bytes4"}],"name":"supportsInterface","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"mint","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"pause","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"unpause","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
const ERC20PresetMinterPauserBin = "0x60806040523480156200001157600080fd5b5060405162001d9538038062001d958339810160408190526200003491620003b5565b8151829082906200004d90600590602085019062000242565b5080516200006390600690602084019062000242565b50506007805460ff19169055506200007d600033620000dd565b620000a97f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a633620000dd565b620000d57f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a33620000dd565b50506200045c565b620000e98282620000ed565b5050565b6200010482826200013060201b6200096e1760201c565b60008281526001602090815260409091206200012b918390620009f2620001d0821b17901c565b505050565b6000828152602081815260408083206001600160a01b038516845290915290205460ff16620000e9576000828152602081815260408083206001600160a01b03851684529091529020805460ff191660011790556200018c3390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b6000620001e7836001600160a01b038416620001f0565b90505b92915050565b60008181526001830160205260408120546200023957508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620001ea565b506000620001ea565b82805462000250906200041f565b90600052602060002090601f016020900481019282620002745760008555620002bf565b82601f106200028f57805160ff1916838001178555620002bf565b82800160010185558215620002bf579182015b82811115620002bf578251825591602001919060010190620002a2565b50620002cd929150620002d1565b5090565b5b80821115620002cd5760008155600101620002d2565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126200031057600080fd5b81516001600160401b03808211156200032d576200032d620002e8565b604051601f8301601f19908116603f01168101908282118183101715620003585762000358620002e8565b816040528381526020925086838588010111156200037557600080fd5b600091505b838210156200039957858201830151818301840152908201906200037a565b83821115620003ab5760008385830101525b9695505050505050565b60008060408385031215620003c957600080fd5b82516001600160401b0380821115620003e157600080fd5b620003ef86838701620002fe565b935060208501519150808211156200040657600080fd5b506200041585828601620002fe565b9150509250929050565b600181811c908216806200043457607f821691505b602082108114156200045657634e487b7160e01b600052602260045260246000fd5b50919050565b611929806200046c6000396000f3fe608060405234801561001057600080fd5b50600436106101c45760003560e01c806370a08231116100f9578063a457c2d711610097578063d539139311610071578063d5391393146103af578063d547741f146103d6578063dd62ed3e146103e9578063e63ab1e91461042257600080fd5b8063a457c2d714610376578063a9059cbb14610389578063ca15c8731461039c57600080fd5b80639010d07c116100d35780639010d07c1461032857806391d148541461035357806395d89b4114610366578063a217fddf1461036e57600080fd5b806370a08231146102e457806379cc67901461030d5780638456cb591461032057600080fd5b8063313ce567116101665780633f4ba83a116101405780633f4ba83a146102ab57806340c10f19146102b357806342966c68146102c65780635c975abb146102d957600080fd5b8063313ce5671461027657806336568abe14610285578063395093511461029857600080fd5b806318160ddd116101a257806318160ddd1461021957806323b872dd1461022b578063248a9ca31461023e5780632f2ff15d1461026157600080fd5b806301ffc9a7146101c957806306fdde03146101f1578063095ea7b314610206575b600080fd5b6101dc6101d73660046115cf565b610449565b60405190151581526020015b60405180910390f35b6101f9610474565b6040516101e89190611625565b6101dc610214366004611674565b610506565b6004545b6040519081526020016101e8565b6101dc61023936600461169e565b61051e565b61021d61024c3660046116da565b60009081526020819052604090206001015490565b61027461026f3660046116f3565b610542565b005b604051601281526020016101e8565b6102746102933660046116f3565b61056d565b6101dc6102a6366004611674565b6105f0565b61027461062f565b6102746102c1366004611674565b6106d5565b6102746102d43660046116da565b610774565b60075460ff166101dc565b61021d6102f236600461171f565b6001600160a01b031660009081526002602052604090205490565b61027461031b366004611674565b610781565b610274610796565b61033b61033636600461173a565b61083a565b6040516001600160a01b0390911681526020016101e8565b6101dc6103613660046116f3565b610859565b6101f9610882565b61021d600081565b6101dc610384366004611674565b610891565b6101dc610397366004611674565b610923565b61021d6103aa3660046116da565b610931565b61021d7f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a681565b6102746103e43660046116f3565b610948565b61021d6103f736600461175c565b6001600160a01b03918216600090815260036020908152604080832093909416825291909152205490565b61021d7f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a81565b60006001600160e01b03198216635a05180f60e01b148061046e575061046e82610a07565b92915050565b60606005805461048390611786565b80601f01602080910402602001604051908101604052809291908181526020018280546104af90611786565b80156104fc5780601f106104d1576101008083540402835291602001916104fc565b820191906000526020600020905b8154815290600101906020018083116104df57829003601f168201915b5050505050905090565b600033610514818585610a3c565b5060019392505050565b60003361052c858285610b60565b610537858585610bf2565b506001949350505050565b60008281526020819052604090206001015461055e8133610dcb565b6105688383610e2f565b505050565b6001600160a01b03811633146105e25760405162461bcd60e51b815260206004820152602f60248201527f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560448201526e103937b632b9903337b91039b2b63360891b60648201526084015b60405180910390fd5b6105ec8282610e51565b5050565b3360008181526003602090815260408083206001600160a01b0387168452909152812054909190610514908290869061062a9087906117d7565b610a3c565b6106597f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a33610859565b6106cb5760405162461bcd60e51b815260206004820152603960248201527f45524332305072657365744d696e7465725061757365723a206d75737420686160448201527f76652070617573657220726f6c6520746f20756e70617573650000000000000060648201526084016105d9565b6106d3610e73565b565b6106ff7f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a633610859565b61076a5760405162461bcd60e51b815260206004820152603660248201527f45524332305072657365744d696e7465725061757365723a206d7573742068616044820152751d99481b5a5b9d195c881c9bdb19481d1bc81b5a5b9d60521b60648201526084016105d9565b6105ec8282610f06565b61077e3382610ff1565b50565b61078c823383610b60565b6105ec8282610ff1565b6107c07f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a33610859565b6108325760405162461bcd60e51b815260206004820152603760248201527f45524332305072657365744d696e7465725061757365723a206d75737420686160448201527f76652070617573657220726f6c6520746f20706175736500000000000000000060648201526084016105d9565b6106d361114b565b600082815260016020526040812061085290836111c6565b9392505050565b6000918252602082815260408084206001600160a01b0393909316845291905290205460ff1690565b60606006805461048390611786565b3360008181526003602090815260408083206001600160a01b0387168452909152812054909190838110156109165760405162461bcd60e51b815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f77604482015264207a65726f60d81b60648201526084016105d9565b6105378286868403610a3c565b600033610514818585610bf2565b600081815260016020526040812061046e906111d2565b6000828152602081905260409020600101546109648133610dcb565b6105688383610e51565b6109788282610859565b6105ec576000828152602081815260408083206001600160a01b03851684529091529020805460ff191660011790556109ae3390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b6000610852836001600160a01b0384166111dc565b60006001600160e01b03198216637965db0b60e01b148061046e57506301ffc9a760e01b6001600160e01b031983161461046e565b6001600160a01b038316610a9e5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b60648201526084016105d9565b6001600160a01b038216610aff5760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b60648201526084016105d9565b6001600160a01b0383811660008181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b6001600160a01b038381166000908152600360209081526040808320938616835292905220546000198114610bec5781811015610bdf5760405162461bcd60e51b815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e636500000060448201526064016105d9565b610bec8484848403610a3c565b50505050565b6001600160a01b038316610c565760405162461bcd60e51b815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f206164604482015264647265737360d81b60648201526084016105d9565b6001600160a01b038216610cb85760405162461bcd60e51b815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201526265737360e81b60648201526084016105d9565b610cc383838361122b565b6001600160a01b03831660009081526002602052604090205481811015610d3b5760405162461bcd60e51b815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e7420657863656564732062604482015265616c616e636560d01b60648201526084016105d9565b6001600160a01b03808516600090815260026020526040808220858503905591851681529081208054849290610d729084906117d7565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610dbe91815260200190565b60405180910390a3610bec565b610dd58282610859565b6105ec57610ded816001600160a01b03166014611236565b610df8836020611236565b604051602001610e099291906117ef565b60408051601f198184030181529082905262461bcd60e51b82526105d991600401611625565b610e39828261096e565b600082815260016020526040902061056890826109f2565b610e5b82826113d2565b60008281526001602052604090206105689082611437565b60075460ff16610ebc5760405162461bcd60e51b815260206004820152601460248201527314185d5cd8589b194e881b9bdd081c185d5cd95960621b60448201526064016105d9565b6007805460ff191690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b6040516001600160a01b03909116815260200160405180910390a1565b6001600160a01b038216610f5c5760405162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f20616464726573730060448201526064016105d9565b610f686000838361122b565b8060046000828254610f7a91906117d7565b90915550506001600160a01b03821660009081526002602052604081208054839290610fa79084906117d7565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b6001600160a01b0382166110515760405162461bcd60e51b815260206004820152602160248201527f45524332303a206275726e2066726f6d20746865207a65726f206164647265736044820152607360f81b60648201526084016105d9565b61105d8260008361122b565b6001600160a01b038216600090815260026020526040902054818110156110d15760405162461bcd60e51b815260206004820152602260248201527f45524332303a206275726e20616d6f756e7420657863656564732062616c616e604482015261636560f01b60648201526084016105d9565b6001600160a01b0383166000908152600260205260408120838303905560048054849290611100908490611864565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b60075460ff16156111915760405162461bcd60e51b815260206004820152601060248201526f14185d5cd8589b194e881c185d5cd95960821b60448201526064016105d9565b6007805460ff191660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258610ee93390565b6000610852838361144c565b600061046e825490565b60008181526001830160205260408120546112235750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561046e565b50600061046e565b610568838383611476565b6060600061124583600261187b565b6112509060026117d7565b67ffffffffffffffff8111156112685761126861189a565b6040519080825280601f01601f191660200182016040528015611292576020820181803683370190505b509050600360fc1b816000815181106112ad576112ad6118b0565b60200101906001600160f81b031916908160001a905350600f60fb1b816001815181106112dc576112dc6118b0565b60200101906001600160f81b031916908160001a905350600061130084600261187b565b61130b9060016117d7565b90505b6001811115611383576f181899199a1a9b1b9c1cb0b131b232b360811b85600f166010811061133f5761133f6118b0565b1a60f81b828281518110611355576113556118b0565b60200101906001600160f81b031916908160001a90535060049490941c9361137c816118c6565b905061130e565b5083156108525760405162461bcd60e51b815260206004820181905260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e7460448201526064016105d9565b6113dc8282610859565b156105ec576000828152602081815260408083206001600160a01b0385168085529252808320805460ff1916905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b6000610852836001600160a01b0384166114dc565b6000826000018281548110611463576114636118b0565b9060005260206000200154905092915050565b60075460ff16156105685760405162461bcd60e51b815260206004820152602a60248201527f45524332305061757361626c653a20746f6b656e207472616e736665722077686044820152691a5b19481c185d5cd95960b21b60648201526084016105d9565b600081815260018301602052604081205480156115c5576000611500600183611864565b855490915060009061151490600190611864565b9050818114611579576000866000018281548110611534576115346118b0565b9060005260206000200154905080876000018481548110611557576115576118b0565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061158a5761158a6118dd565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061046e565b600091505061046e565b6000602082840312156115e157600080fd5b81356001600160e01b03198116811461085257600080fd5b60005b838110156116145781810151838201526020016115fc565b83811115610bec5750506000910152565b60208152600082518060208401526116448160408501602087016115f9565b601f01601f19169190910160400192915050565b80356001600160a01b038116811461166f57600080fd5b919050565b6000806040838503121561168757600080fd5b61169083611658565b946020939093013593505050565b6000806000606084860312156116b357600080fd5b6116bc84611658565b92506116ca60208501611658565b9150604084013590509250925092565b6000602082840312156116ec57600080fd5b5035919050565b6000806040838503121561170657600080fd5b8235915061171660208401611658565b90509250929050565b60006020828403121561173157600080fd5b61085282611658565b6000806040838503121561174d57600080fd5b50508035926020909101359150565b6000806040838503121561176f57600080fd5b61177883611658565b915061171660208401611658565b600181811c9082168061179a57607f821691505b602082108114156117bb57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b600082198211156117ea576117ea6117c1565b500190565b7f416363657373436f6e74726f6c3a206163636f756e74200000000000000000008152600083516118278160178501602088016115f9565b7001034b99036b4b9b9b4b733903937b6329607d1b60179184019182015283516118588160288401602088016115f9565b01602801949350505050565b600082821015611876576118766117c1565b500390565b6000816000190483118215151615611895576118956117c1565b500290565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b6000816118d5576118d56117c1565b506000190190565b634e487b7160e01b600052603160045260246000fdfea2646970667358221220bae3c8fa2cbf779ea7b42e35284a76579a8a5958ffcec2fe816eda805f8e039d64736f6c634300080b0033"
